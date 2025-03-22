package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yeyus/gumble-joselito/pkg/dmr"
)

type MessageGroupJoin struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type   MessageType
	Groups []*dmr.DMRID
}

func NewMessageGroupJoin(groups []*dmr.DMRID) *MessageGroupJoin {
	return &MessageGroupJoin{
		Type:   GROUP_JOIN,
		Groups: groups,
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
