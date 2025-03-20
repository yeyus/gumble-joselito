package joselito

import (
	"reflect"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestDMRIDMarshall_FixInt(t *testing.T) {
	id := &DMRID{Id: uint32(93)}

	b, err := msgpack.Marshal(id)
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0x5D}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestDMRIDMarshall_Uint8(t *testing.T) {
	id := &DMRID{Id: uint32(214)}

	b, err := msgpack.Marshal(id)
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0xCC, 0xD6}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestDMRIDMarshall_Uint16(t *testing.T) {
	id := &DMRID{Id: uint32(3112)}

	b, err := msgpack.Marshal(id)
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0xCD, 0x0C, 0x28}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}

func TestDMRIDMarshall_Uint32(t *testing.T) {
	id := &DMRID{Id: uint32(520312)}

	b, err := msgpack.Marshal(id)
	if err != nil {
		t.Errorf("marshalling returned error: %v", err)
	}

	expected := []byte{0xCE, 0x00, 0x07, 0xF0, 0x78}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("payload doesn't match, expected %v received %v", expected, b)
	}
}
