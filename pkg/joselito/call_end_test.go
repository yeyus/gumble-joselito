package joselito

import (
	"reflect"
	"testing"
)

func TestMessageCallEndMarshall(t *testing.T) {
	message := NewMessageCallEnd()

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x91, 0x0D}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageCallEndUnmarshall(t *testing.T) {
	entity := NewMessageCallEnd()

	msg := []byte{0x91, 0x0D}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != CALL_END {
		t.Errorf("unmarshalling expected type to be CALL_DROP, received %d", entity.Type)
	}
}
