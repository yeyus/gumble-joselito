package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageCallDrop struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type MessageType
}

func NewMessageCallDrop() *MessageCallDrop {
	return &MessageCallDrop{
		Type: CALL_DROP,
	}
}

func (m *MessageCallDrop) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallDrop) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallDrop) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_DROP {
		return errors.New("type mismatch")
	}
	return err
}
