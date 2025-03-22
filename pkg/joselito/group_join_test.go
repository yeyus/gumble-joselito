package joselito

import (
	"reflect"
	"testing"

	"github.com/yeyus/gumble-joselito/pkg/dmr"
)

func TestMessageGroupJoinMarshall_SingleGroup(t *testing.T) {
	message := NewMessageGroupJoin([]*dmr.DMRID{dmr.NewDMRID(214)})

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x92, 0x01, 0x91, 0xCC, 0xD6}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageGroupJoinUnmarshall_SingleGroup(t *testing.T) {
	entity := NewMessageGroupJoin(nil)

	msg := []byte{0x92, 0x01, 0x91, 0xCC, 0xD6}
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

func TestMessageGroupJoinMarshall_MultipleGroup(t *testing.T) {
	message := NewMessageGroupJoin([]*dmr.DMRID{dmr.NewDMRID(93), dmr.NewDMRID(222)})

	b, err := message.Marshall()
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x92, 0x01, 0x92, 0x5D, 0xCC, 0xDE}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestMessageGroupJoinUnmarshall_MultipleGroup(t *testing.T) {
	entity := NewMessageGroupJoin(nil)

	msg := []byte{0x92, 0x01, 0x92, 0xCE, 0x00, 0x07, 0xF0, 0x78, 0x5D}
	err := entity.Unmarshall(msg)
	if err != nil {
		t.Errorf("unmarshalling returned error: %v", err)
	}

	if len(entity.Groups) != 2 {
		t.Errorf("unmarshalling expected 2 groups, received %d", len(entity.Groups))
	}

	if entity.Groups[0].Id != 520312 {
		t.Errorf("unmarshalling values failed id expected to be 520312, received %d", entity.Groups[0].Id)
	}

	if entity.Groups[1].Id != 93 {
		t.Errorf("unmarshalling values failed id expected to be 93, received %d", entity.Groups[1].Id)
	}
}
