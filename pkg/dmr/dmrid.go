package dmr

import (
	"fmt"
	"strconv"

	"github.com/enescakir/emoji"
	"github.com/vmihailenco/msgpack/v5"
)

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

func (d *DMRID) StringWithEmoji() string {
	flag, err := d.CountryEmoji()
	if err != nil {
		return fmt.Sprintf("%d", d.Id)
	}

	return fmt.Sprintf("%v%d", flag, d.Id)
}

func (d *DMRID) ISOCountry() (string, error) {
	if d.Id < 100 {
		return "wo", nil
	}

	prefix := fmt.Sprintf("%d", d.Id)[:3]

	num, err := strconv.ParseInt(prefix, 10, 32)
	if err != nil {
		return "", err
	}

	val, ok := countryMap[num]
	if ok {
		return val, nil
	}

	return "wo", nil
}

func (d *DMRID) CountryEmoji() (emoji.Emoji, error) {
	country, err := d.ISOCountry()
	if err != nil {
		return emoji.QuestionMark, err

	}

	if country == "wo" {
		return emoji.GlobeShowingEuropeAfrica, nil
	}

	return emoji.CountryFlag(country)
}

var _ msgpack.CustomEncoder = (*DMRID)(nil)

var _ msgpack.CustomDecoder = (*DMRID)(nil)
