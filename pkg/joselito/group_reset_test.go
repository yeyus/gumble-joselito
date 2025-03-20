package joselito

import (
	"reflect"
	"testing"
)

func TestMessageGroupResetMarshall(t *testing.T) {
	message := NewMessageGroupReset()

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x91, 0x03}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageGroupResetUnmarshall(t *testing.T) {
	entity := NewMessageGroupReset()

	msg := []byte{0x91, 0x03}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != GROUP_RESET {
		t.Errorf("unmarshalling expected message type to be GROUP_RESET, received %d", entity.Type)
	}
}
