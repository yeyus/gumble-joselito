package audio

import "testing"

func TestJLaw(t *testing.T) {
	tests := map[uint8]int16{0: -32124, 127: 0, 130: 30076, 255: 0, 15: -16764}

	for i, o := range tests {
		if ULawDecode[i] != o {
			t.Errorf("compansion values don't match input %d should produce %d instead produced %d", i, o, ULawDecode[i])
		}
	}

}
