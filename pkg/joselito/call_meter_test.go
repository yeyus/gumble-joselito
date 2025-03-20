package joselito

import (
	"reflect"
	"testing"
)

func TestMessageCallMeterMarshall(t *testing.T) {
	message := NewMessageCallMeter(-22.891209)

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x92, 0x16, 0xCA, 0xC1, 0xB7, 0x21, 0x32}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageCallMeterUnmarshall(t *testing.T) {
	entity := NewMessageCallMeter(0)

	msg := []byte{0x92, 0x16, 0xCA, 0xC1, 0xB7, 0x21, 0x32}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != CALL_METER {
		t.Errorf("unmarshalling expected type to be CALL_METER, received %d", entity.Type)
	}

	expected := float32(-22.891209)
	if entity.Volume != expected {
		t.Errorf("volume doesn't match, expected \"%f\" received %f", expected, entity.Volume)
	}
}
