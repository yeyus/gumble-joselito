package joselito

import (
	"reflect"
	"testing"
)

func TestMessageCallDropMarshall(t *testing.T) {
	message := NewMessageCallDrop()

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x91, 0x0C}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageCallDropUnmarshall(t *testing.T) {
	entity := NewMessageCallDrop()

	msg := []byte{0x91, 0x0C}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != CALL_DROP {
		t.Errorf("unmarshalling expected type to be CALL_DROP, received %d", entity.Type)
	}
}
