package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageCallMeter struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type   MessageType
	Volume float32
}

func NewMessageCallMeter(volume float32) *MessageCallMeter {
	return &MessageCallMeter{
		Type:   CALL_METER,
		Volume: volume,
	}
}

func (m *MessageCallMeter) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallMeter) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallMeter) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_METER {
		return errors.New("type mismatch")
	}
	return err
}
