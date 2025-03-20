package joselito

import "github.com/vmihailenco/msgpack/v5"

type DMRID struct {
	Id uint
}

func NewDMRID(id uint) *DMRID {
	return &DMRID{Id: id}
}

// DecodeMsgpack implements msgpack.CustomDecoder.
func (d *DMRID) DecodeMsgpack(dec *msgpack.Decoder) error {
	s, err := dec.DecodeUint()
	d.Id = s
	return err
}

// EncodeMsgpack implements msgpack.CustomEncoder.
func (d *DMRID) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeUint(uint64(d.Id))
}

var _ msgpack.CustomEncoder = (*DMRID)(nil)

var _ msgpack.CustomDecoder = (*DMRID)(nil)
