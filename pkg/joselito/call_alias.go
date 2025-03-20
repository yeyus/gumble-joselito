package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageCallAlias struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type        MessageType
	TalkerAlias string
}

func NewMessageCallAlias(talkerAlias string) *MessageCallAlias {
	return &MessageCallAlias{
		Type:        CALL_ALIAS,
		TalkerAlias: talkerAlias,
	}
}

func (m *MessageCallAlias) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallAlias) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallAlias) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_ALIAS {
		return errors.New("type mismatch")
	}
	return err
}
