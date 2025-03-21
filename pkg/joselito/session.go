package joselito

import (
	"errors"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type OnCallStartCallback func(*Call, *MessageCallStart) error
type OnCallEndCallback func(*Call, *MessageCallEnd) error
type OnCallMeterUpdateCallback func(*Call, *MessageCallMeter) error
type OnCallTalkerAliasUpdateCallback func(*Call, *MessageCallAlias) error
type OnCallAudioReceivedCallback func(*Call, *MessageCallAudio) error

type Session struct {
	logger *log.Logger

	connection *websocket.Conn
	Talkgroups []*DMRID

	// Last Call
	Call *Call

	callStartCallbacks             []OnCallStartCallback
	callEndCallbacks               []OnCallEndCallback
	callMeterUpdateCallbacks       []OnCallMeterUpdateCallback
	callTalkerAliasUpdateCallbacks []OnCallTalkerAliasUpdateCallback
	callAudioReceivedCallbacks     []OnCallAudioReceivedCallback

	SessionEnd chan struct{}
}

func NewSession(connection *websocket.Conn) *Session {
	session := &Session{
		logger:                         log.New(os.Stdout, "[joselito-session] ", log.LstdFlags),
		connection:                     connection,
		Talkgroups:                     make([]*DMRID, 0),
		SessionEnd:                     make(chan struct{}),
		callStartCallbacks:             make([]OnCallStartCallback, 0),
		callEndCallbacks:               make([]OnCallEndCallback, 0),
		callMeterUpdateCallbacks:       make([]OnCallMeterUpdateCallback, 0),
		callTalkerAliasUpdateCallbacks: make([]OnCallTalkerAliasUpdateCallback, 0),
		callAudioReceivedCallbacks:     make([]OnCallAudioReceivedCallback, 0),
	}

	go func() {
		defer close(session.SessionEnd)
		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				session.logger.Println("error reading from websocket:", err)
				break
			}

			err = session.ProcessMessage(messageType, message)
			if err != nil {
				session.logger.Printf("error processing protocol message: %v", err)
			}

			// session.logger.Printf("recv: %s", message)
		}
		session.logger.Printf("websocket loop died")
	}()

	return session
}

func (s *Session) GroupJoin(talkgroups []*DMRID) error {
	s.Talkgroups = append(s.Talkgroups, talkgroups...)
	joinMessage, err := NewMessageGroupJoin(s.Talkgroups).Marshall()
	if err != nil {
		return err
	}

	return s.connection.WriteMessage(websocket.BinaryMessage, joinMessage)
}

func (s *Session) ProcessMessage(messageType int, buffer []byte) error {
	if messageType != websocket.BinaryMessage {
		return nil
	}

	callAliasMsg := NewMessageCallAlias("")
	err := callAliasMsg.Unmarshall(buffer)
	if err == nil {
		return s.onCallAlias(callAliasMsg)
	}

	callAudioMsg := NewMessageCallAudio(nil)
	err = callAudioMsg.Unmarshall(buffer)
	if err == nil {
		return s.onCallAudio(callAudioMsg)
	}

	callStartMsg := NewMessageCallStart(nil, nil)
	err = callStartMsg.Unmarshall(buffer)
	if err == nil {
		return s.onCallStart(callStartMsg)
	}

	callEndMsg := NewMessageCallEnd()
	err = callEndMsg.Unmarshall(buffer)
	if err == nil {
		return s.onCallEnd(callEndMsg)
	}

	callMeterMsg := NewMessageCallMeter(0)
	err = callMeterMsg.Unmarshall(buffer)
	if err == nil {
		return s.onCallMeter(callMeterMsg)
	}

	// unknown message type
	s.logger.Printf("unknown message received: %v", buffer)
	return errors.New("unknown message type")
}

func (s *Session) AddOnCallAliasCallback(cb OnCallTalkerAliasUpdateCallback) {
	s.callTalkerAliasUpdateCallbacks = append(s.callTalkerAliasUpdateCallbacks, cb)
}

func (s *Session) onCallAlias(msg *MessageCallAlias) error {
	s.logger.Printf("call alias: %s", msg.TalkerAlias)

	if s.Call != nil {
		s.Call.TalkerAlias = msg.TalkerAlias
	}

	for _, cb := range s.callTalkerAliasUpdateCallbacks {
		err := cb(s.Call, msg)
		if err != nil {
			s.logger.Printf("callTalkerAliasUpdateCallbacks err: %v", err)
		}
	}

	return nil
}

func (s *Session) AddOnCallAudioReceivedCallback(cb OnCallAudioReceivedCallback) {
	s.callAudioReceivedCallbacks = append(s.callAudioReceivedCallbacks, cb)
}

func (s *Session) onCallAudio(msg *MessageCallAudio) error {
	s.logger.Printf("call audio received: %d", len(msg.Data))

	for _, cb := range s.callAudioReceivedCallbacks {
		err := cb(s.Call, msg)
		if err != nil {
			s.logger.Printf("callAudioReceivedCallbacks err: %v", err)
		}
	}

	return nil
}

func (s *Session) AddOnCallStartCallback(cb OnCallStartCallback) {
	s.callStartCallbacks = append(s.callStartCallbacks, cb)
}

func (s *Session) onCallStart(msg *MessageCallStart) error {
	s.Call = NewCall(msg.Origin, msg.Destination)

	s.logger.Printf("%s call start", s.Call)

	for _, cb := range s.callStartCallbacks {
		err := cb(s.Call, msg)
		if err != nil {
			s.logger.Printf("callStartCallbacks err: %v", err)
		}
	}

	return nil
}

func (s *Session) AddOnCallEndCallback(cb OnCallEndCallback) {
	s.callEndCallbacks = append(s.callEndCallbacks, cb)
}

func (s *Session) onCallEnd(msg *MessageCallEnd) error {

	if s.Call != nil {
		s.Call.Finish()
	}

	for _, cb := range s.callEndCallbacks {
		err := cb(s.Call, msg)
		if err != nil {
			s.logger.Printf("callEndCallbacks err: %v", err)
		}
	}

	s.logger.Printf("%s call end", s.Call)

	return nil
}

func (s *Session) AddOnCallMeterUpdateCallback(cb OnCallMeterUpdateCallback) {
	s.callMeterUpdateCallbacks = append(s.callMeterUpdateCallbacks, cb)
}

func (s *Session) onCallMeter(msg *MessageCallMeter) error {
	// s.logger.Printf("call meter: %f", msg.Volume)

	if s.Call != nil {
		s.Call.Volume = msg.Volume
	}

	for _, cb := range s.callMeterUpdateCallbacks {
		err := cb(s.Call, msg)
		if err != nil {
			s.logger.Printf("callMeterUpdateCallbacks err: %v", err)
		}
	}

	return nil
}
