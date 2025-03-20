package joselito

import (
	"reflect"
	"testing"
)

func TestMessageCallAliasMarshall(t *testing.T) {
	message := NewMessageCallAlias("EA7JMF - Jesus")

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x92, 0x15, 0xAE, 'E', 'A', '7', 'J', 'M', 'F', ' ', '-', ' ', 'J', 'e', 's', 'u', 's'}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageCallAliasUnmarshall(t *testing.T) {
	entity := NewMessageCallAlias("")

	msg := []byte{0x92, 0x15, 0xAE, 'E', 'A', '7', 'J', 'M', 'F', ' ', '-', ' ', 'J', 'e', 's', 'u', 's'}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if entity.Type != CALL_ALIAS {
		t.Errorf("unmarshalling expected type to be CALL_ALIAS, received %d", entity.Type)
	}

	expected := "EA7JMF - Jesus"
	if entity.TalkerAlias != expected {
		t.Errorf("talker alias doesn't match, expected \"%s\" received %s", expected, entity.TalkerAlias)
	}
}
