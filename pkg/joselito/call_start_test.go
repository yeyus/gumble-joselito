package joselito

import (
	"reflect"
	"testing"

	"github.com/yeyus/gumble-joselito/pkg/dmr"
)

func TestMessageCallStartMarshall(t *testing.T) {
	message := NewMessageCallStart(dmr.NewDMRID(3101456), dmr.NewDMRID(93))

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x95, 0x0B, 0x00, 0xCE, 0x00, 0x2F, 0x53, 0x10, 0x5D, 0x00}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageCallStartUnmarshall(t *testing.T) {
	entity := NewMessageCallStart(nil, nil)

	msg := []byte{0x95, 0x0B, 0x00, 0xCE, 0x00, 0x2F, 0x53, 0x10, 0x5D, 0x00}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != CALL_START {
		t.Errorf("unmarshalling expected message type to be CALL_START, received %d", entity.Type)
	}

	if entity.Origin.Id != 3101456 {
		t.Errorf("unmarshalling expected origin to be 3101456, received %d", entity.Origin.Id)
	}

	if entity.Destination.Id != 93 {
		t.Errorf("unmarshalling expected destination to be 93, received %d", entity.Destination.Id)
	}
}
