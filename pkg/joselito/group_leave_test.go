package joselito

import (
	"reflect"
	"testing"
)

func TestMessageGroupLeaveMarshall_SingleGroup(t *testing.T) {
	message := NewMessageGroupLeave()
	message.Groups = []*DMRID{NewDMRID(214)}

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x92, 0x02, 0x91, 0xCC, 0xD6}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageGroupLeaveUnmarshall_SingleGroup(t *testing.T) {
	entity := NewMessageGroupLeave()

	msg := []byte{0x92, 0x02, 0x91, 0xCC, 0xD6}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if len(entity.Groups) != 1 {
		t.Errorf("unmarshalling expected 1 group, received %d", len(entity.Groups))
	}

	if entity.Groups[0].Id != 214 {
		t.Errorf("unmarshalling values failed id expected to be 214, received %d", entity.Groups[0].Id)
	}
}
