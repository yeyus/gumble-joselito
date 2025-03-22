package joselito

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yeyus/gumble-joselito/pkg/audio"
	"github.com/yeyus/gumble/gumble"
)

// State represents the state of a Stream.
type State int32

// Valid states of Stream.
const (
	StateInitial State = iota + 1
	StatePlaying
	StatePaused
	StateStopped
)

type Stream struct {
	logger *log.Logger

	session *Session
	state   State
	client  *gumble.Client
	pause   chan struct{}

	pipeWriter *io.PipeWriter
	pipeReader *io.PipeReader

	// Playback volume (can be changed while the source is playing).
	Volume float32

	// per struct buffers
	pcm       []int16
	upsampled *bytes.Buffer

	lock sync.Mutex
}

func NewStream(client *gumble.Client, session *Session) *Stream {
	r, w := io.Pipe()

	stream := &Stream{
		logger:     log.New(os.Stdout, "[joselito-stream] ", log.LstdFlags),
		client:     client,
		session:    session,
		Volume:     1.0,
		pause:      make(chan struct{}),
		state:      StateInitial,
		pipeWriter: w,
		pipeReader: r,

		// init buffers
		pcm:       make([]int16, 960),
		upsampled: new(bytes.Buffer),
	}

	session.AddOnCallStartCallback(stream.onCallStart)
	session.AddOnCallEndCallback(stream.onCallEnd)
	session.AddOnCallAudioReceivedCallback(stream.onCallAudio)

	return stream
}

// Audio rate is 8000Hz, 8bit uLaw
// 8bit * 960samples at 8.33Hz ws rate
// conversion 960 (8bitsamples) * (interpolation factor M=6) * 2 channels= 11520 samples
var UPSAMPLE_FACTOR = 6

func (s *Stream) onCallAudio(call *Call, msg *MessageCallAudio) error {
	// do the conversion here

	for i := range msg.Data {
		s.pcm[i] = audio.ULawDecode[msg.Data[i]]
	}

	// should output 960*6 in int16, len 5670 16 bit samples
	upsampledAndFiltered := audio.LinearUpsampler(s.pcm, UPSAMPLE_FACTOR)

	// upsample by factor M=6, len upsampled*2bytes(int16)
	for i := 0; i < len(upsampledAndFiltered); i += 1 {
		err := binary.Write(s.upsampled, binary.LittleEndian, upsampledAndFiltered[i])
		if err != nil {
			s.logger.Fatal(err)
		}
	}

	_, err := s.upsampled.WriteTo(s.pipeWriter)
	if err != nil {
		s.logger.Printf("onCallAudio: problem writing audio samples: %v", err)
		return err
	}
	s.upsampled.Reset()

	return nil
}

func (s *Stream) onCallStart(call *Call, msg *MessageCallStart) error {
	s.logger.Printf("call started %s", call.String())
	return s.Play()
}

func (s *Stream) onCallEnd(call *Call, msg *MessageCallEnd) error {
	s.logger.Printf("call end %s", call.String())
	go s.Pause()
	return nil
}

func (s *Stream) Play() error {
	s.logger.Printf("starting stream")
	s.lock.Lock()
	defer s.lock.Unlock()

	switch s.state {
	case StatePaused:
		s.state = StatePlaying
		go s.process()
		return nil
	case StatePlaying:
		return errors.New("stream already playing")
	case StateStopped:
		return errors.New("stream has stopped")
	}

	s.state = StatePlaying
	go s.process()
	return nil
}

func (s *Stream) Pause() error {
	s.logger.Printf("pausing stream")
	s.lock.Lock()
	if s.state != StatePlaying {
		s.lock.Unlock()
		return errors.New("stream is not playing")
	}
	s.state = StatePaused
	s.lock.Unlock()
	s.pause <- struct{}{}

	return nil
}

func (s *Stream) Stop() error {
	s.lock.Lock()
	switch s.state {
	case StateStopped, StateInitial:
		s.lock.Unlock()
		return errors.New("stream is not playing nor paused")
	}
	s.cleanup()

	return nil
}

// State returns the state of the stream.
func (s *Stream) State() State {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.state
}

func (s *Stream) process() {
	interval := s.client.Config.AudioInterval
	frameSize := s.client.Config.AudioFrameSize()

	byteBuffer := make([]byte, frameSize*2)

	outgoing := s.client.AudioOutgoing()
	defer close(outgoing)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.pause:
			s.logger.Printf("stream paused")
			return
		case <-ticker.C:
			if _, err := io.ReadFull(s.pipeReader, byteBuffer); err != nil {
				s.lock.Lock()
				s.cleanup()
				return
			}
			int16Buffer := make([]int16, frameSize)
			for i := range int16Buffer {
				float := float32(int16(binary.LittleEndian.Uint16(byteBuffer[i*2 : (i+1)*2])))
				int16Buffer[i] = int16(s.Volume * float)
			}
			outgoing <- gumble.AudioBuffer(int16Buffer)
		}
	}
}

func (s *Stream) cleanup() {
	defer s.lock.Unlock()
	// s.l has been acquired
	if s.state == StateStopped {
		return
	}

	for len(s.pause) > 0 {
		<-s.pause
	}
	s.state = StateStopped
}
