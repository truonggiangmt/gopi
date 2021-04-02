package chromecast

import (
	context "context"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	empty "github.com/golang/protobuf/ptypes/empty"
)

type Service struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Server
	gopi.CastManager
	gopi.Publisher
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *Service) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.Server, this.CastManager, this.Publisher)

	if err := this.RegisterService(RegisterManagerServer, this); err != nil {
		return err
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// RPC METHODS

func (this *Service) List(context.Context, *empty.Empty) (*ListResponse, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Service) Stream(_ *empty.Empty, stream Manager_StreamServer) error {
	this.Logger.Debug("<Stream>")

	// Send a null event once a second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Subscribe to input events
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Obtain server cancel context
	ctx := this.Server.NewStreamContext()

	// Loop which streams until server context cancels or an error occurs sending a Ping
	for {
		select {
		case evt := <-ch:
			// Send event
			if evt, ok := evt.(gopi.CastEvent); ok {
				this.Debug("Stream: ", evt)
				if err := stream.Send(toProtoEvent(evt)); err != nil {
					this.Print("Stream: ", err)
				}
			}
		case <-ctx.Done():
			// Context done
			return ctx.Err()
		case <-ticker.C:
			// Send a ping
			if err := stream.Send(toProtoNull()); err != nil {
				this.Debug("Stream: ", "Error sending null event, ending stream")
				return err
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Service) mustEmbedUnimplementedManagerServer() {
	// NOOP
}
