package dmr

import (
	"fmt"
	"time"
)

type Call struct {
	Origin      *DMRID
	Destination *DMRID

	TalkerAlias string
	Volume      float32

	Start time.Time
	End   time.Time
}

func NewCall(origin *DMRID, destination *DMRID) *Call {
	return &Call{
		Origin:      origin,
		Destination: destination,
		TalkerAlias: "",
		Volume:      0,
		Start:       time.Now(),
	}
}

func (c *Call) Finished() bool {
	return !c.End.IsZero()
}

func (c *Call) Finish() {
	c.End = time.Now()
}

func (c *Call) Duration() time.Duration {
	end := c.End
	if end.IsZero() {
		end = time.Now()
	}

	return end.Sub(c.Start)
}

func (c *Call) String() string {
	return fmt.Sprintf("[Call de:%v to:%v duration:%s talkeralias=%s]", c.Origin.StringWithEmoji(), c.Destination.StringWithEmoji(), c.Duration().String(), c.TalkerAlias)
}
