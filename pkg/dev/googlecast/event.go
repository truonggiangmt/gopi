package googlecast

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	cast  *Cast
	flags gopi.CastFlag
}

type state struct {
	key    string
	req    int
	values []interface{}
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(device *Cast, flags gopi.CastFlag) gopi.CastEvent {
	return &event{device, flags}
}

func NewState(key string, req int, values ...interface{}) state {
	return state{key, req, values}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *event) Name() string {
	if this.cast != nil {
		return this.cast.Name()
	} else {
		return ""
	}
}

func (this *event) Cast() gopi.Cast {
	return this.cast
}

func (this *event) Flags() gopi.CastFlag {
	return this.flags
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<cast.event"
	if this.flags != gopi.CAST_FLAG_NONE {
		str += " flags=" + fmt.Sprint(this.flags)
	}
	if this.cast != nil {
		str += " device=" + fmt.Sprint(this.cast)
	}
	return str + ">"
}