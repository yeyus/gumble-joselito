package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type MessageCallAudio struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type MessageType
	Data []byte
}

func NewMessageCallAudio(buffer []byte) *MessageCallAudio {
	return &MessageCallAudio{
		Type: CALL_AUDIO,
		Data: buffer,
	}
}

func (m *MessageCallAudio) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallAudio) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallAudio) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_AUDIO {
		return errors.New("type mismatch")
	}
	return err
}
