package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageCallEnd struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type MessageType
}

func NewMessageCallEnd() *MessageCallEnd {
	return &MessageCallEnd{
		Type: CALL_END,
	}
}

func (m *MessageCallEnd) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallEnd) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallEnd) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_END {
		return errors.New("type mismatch")
	}
	return err
}
