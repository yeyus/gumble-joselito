package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageGroupReset struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type MessageType
}

func NewMessageGroupReset() *MessageGroupReset {
	return &MessageGroupReset{
		Type: GROUP_RESET,
	}
}

func (m *MessageGroupReset) MessageType() MessageType {
	return m.Type
}

func (m *MessageGroupReset) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageGroupReset) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != GROUP_RESET {
		return errors.New("type mismatch")
	}
	return err
}
