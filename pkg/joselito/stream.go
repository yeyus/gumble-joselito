package joselito

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

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
	}

	session.AddOnCallStartCallback(stream.onCallStart)
	session.AddOnCallEndCallback(stream.onCallEnd)
	session.AddOnCallAudioReceivedCallback(stream.onCallAudio)

	return stream
}

func (s *Stream) onCallAudio(call *Call, msg *MessageCallAudio) error {
	_, err := s.pipeWriter.Write(msg.Data)
	if err != nil {
		s.logger.Printf("onCallAudio: problem writing audio samples: %v", err)
		return err
	}

	return nil
}

func (s *Stream) onCallStart(call *Call, msg *MessageCallStart) error {
	s.logger.Printf("call started %s", call.String())
	return s.Play()
}

func (s *Stream) onCallEnd(call *Call, msg *MessageCallEnd) error {
	s.logger.Printf("call end %s", call.String())
	return s.Pause()
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
	s.logger.Printf("stream paused")
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
			return
		case <-ticker.C:
			if _, err := io.ReadFull(s.pipeReader, byteBuffer); err != nil {
				s.lock.Lock()
				s.cleanup()
				return
			}
			int16Buffer := make([]int16, frameSize)
			for i := range int16Buffer {
				float := float32(int16(binary.BigEndian.Uint16(byteBuffer[i*2 : (i+1)*2])))
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
