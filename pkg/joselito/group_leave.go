package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageGroupLeave struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type   MessageType
	Groups []*DMRID
}

func NewMessageGroupLeave() *MessageGroupLeave {
	return &MessageGroupLeave{
		Type: GROUP_LEAVE,
	}
}

func (m *MessageGroupLeave) MessageType() MessageType {
	return m.Type
}

func (m *MessageGroupLeave) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageGroupLeave) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != GROUP_LEAVE {
		return errors.New("type mismatch")
	}
	return err
}
