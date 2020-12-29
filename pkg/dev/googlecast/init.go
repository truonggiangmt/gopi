package googlecast

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	graph "github.com/djthorpe/gopi/v3/pkg/graph"
)

func init() {
	// Register gopi.CastManager
	graph.RegisterUnit(reflect.TypeOf(&Manager{}), reflect.TypeOf((*gopi.CastManager)(nil)))
}
