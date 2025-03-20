package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageGroupJoin struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type   MessageType
	Groups []*DMRID
}

func NewMessageGroupJoin() *MessageGroupJoin {
	return &MessageGroupJoin{
		Type: GROUP_JOIN,
	}
}

func (m *MessageGroupJoin) MessageType() MessageType {
	return m.Type
}

func (m *MessageGroupJoin) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageGroupJoin) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != GROUP_JOIN {
		return errors.New("type mismatch")
	}
	return err
}
