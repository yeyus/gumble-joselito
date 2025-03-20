package joselito

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

type Session struct {
	connection *websocket.Conn
	Talkgroups []*DMRID

	// current call

	SessionEnd chan struct{}
}

func NewSession(connection *websocket.Conn) *Session {
	session := &Session{
		connection: connection,
		Talkgroups: make([]*DMRID, 0),
		SessionEnd: make(chan struct{}),
	}

	go func() {
		defer close(session.SessionEnd)
		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				log.Println("error reading from websocket:", err)
				return
			}

			err = session.ProcessMessage(messageType, message)
			if err != nil {
				log.Printf("error processing protocol message: %v", err)
			}

			// log.Printf("recv: %s", message)
		}
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
	log.Printf("unknown message received: %v", buffer)
	return errors.New("unknown message type")
}

func (s *Session) onCallAlias(msg *MessageCallAlias) error {
	log.Printf("call alias: %s", msg.TalkerAlias)
	return nil
}

func (s *Session) onCallAudio(msg *MessageCallAudio) error {
	log.Printf("call audio received: %d", len(msg.Data))
	return nil
}

func (s *Session) onCallStart(msg *MessageCallStart) error {
	log.Printf("call start: from %d to %d", msg.Origin.Id, msg.Destination.Id)
	return nil
}

func (s *Session) onCallEnd(msg *MessageCallEnd) error {
	log.Printf("call end")
	return nil
}

func (s *Session) onCallMeter(msg *MessageCallMeter) error {
	log.Printf("call meter: %f", msg.Volume)
	return nil
}
