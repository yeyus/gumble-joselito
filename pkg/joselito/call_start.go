package joselito

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yeyus/gumble-joselito/pkg/dmr"
)

type MessageCallStart struct {
	_msgpack struct{} `msgpack:",as_array"`

	Type        MessageType
	Unknown1    uint
	Origin      *dmr.DMRID
	Destination *dmr.DMRID
	Unknown2    uint
}

func NewMessageCallStart(origin *dmr.DMRID, destination *dmr.DMRID) *MessageCallStart {
	return &MessageCallStart{
		Type:        CALL_START,
		Unknown1:    0,
		Origin:      origin,
		Destination: destination,
		Unknown2:    0,
	}
}

func (m *MessageCallStart) MessageType() MessageType {
	return m.Type
}

func (m *MessageCallStart) Marshall() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *MessageCallStart) Unmarshall(buffer []byte) error {
	err := msgpack.Unmarshal(buffer, m)

	if m.Type != CALL_START {
		return errors.New("type mismatch")
	}
	return err
}
