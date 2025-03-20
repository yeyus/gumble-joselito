package joselito

import "github.com/vmihailenco/msgpack/v5"

type DMRID struct {
	Id uint32
}

func NewDMRID(id uint32) *DMRID {
	return &DMRID{Id: id}
}

// DecodeMsgpack implements msgpack.CustomDecoder.
func (d *DMRID) DecodeMsgpack(dec *msgpack.Decoder) error {
	s, err := dec.DecodeUint8()
	if err == nil {
		d.Id = uint32(s)
		return nil
	}

	t, err := dec.DecodeUint16()
	if err == nil {
		d.Id = uint32(t)
		return nil
	}

	u, err := dec.DecodeUint32()
	if err != nil {
		return err
	}

	d.Id = u
	return nil
}

// EncodeMsgpack implements msgpack.CustomEncoder.
func (d *DMRID) EncodeMsgpack(enc *msgpack.Encoder) error {
	if d.Id <= 0xFF {
		return enc.EncodeUint8(uint8(d.Id))
	} else if d.Id <= 0xFFFF {
		return enc.EncodeUint16(uint16(d.Id))
	}

	return enc.EncodeUint32(d.Id)
}

var _ msgpack.CustomEncoder = (*DMRID)(nil)

var _ msgpack.CustomDecoder = (*DMRID)(nil)
