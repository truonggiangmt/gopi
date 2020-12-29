package main

import (
	"context"
	"net"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.ConnPool
	gopi.Command
	gopi.Logger
	gopi.PingService
	gopi.MetricsService
	gopi.Server
	gopi.ServiceDiscovery
	gopi.Unit
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("server", "Start ping service", func(ctx context.Context) error {
		if network, addr, err := this.GetServeAddress(); err != nil {
			return err
		} else {
			return this.RunServer(ctx, network, addr)
		}
	})
	cfg.Command("version", "Display server version information", func(ctx context.Context) error {
		if stub, err := this.GetStub(); err != nil {
			return err
		} else {
			return this.RunVersion(ctx, stub)
		}
	})
	cfg.Command("ping", "Perform ping to server", func(ctx context.Context) error {
		if stub, err := this.GetStub(); err != nil {
			return err
		} else {
			return this.RunPing(ctx, stub)
		}
	})
	cfg.Command("metrics", "Retrieve metrics from server", func(ctx context.Context) error {
		if stub, err := this.GetMetricsStub(); err != nil {
			return err
		} else {
			return this.RunMetrics(ctx, stub)
		}
	})

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if cmd, err := cfg.GetCommand(cfg.Args()); err != nil {
		return gopi.ErrHelp
	} else if cmd == nil {
		return gopi.ErrHelp
	} else {
		this.Command = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) GetStub() (gopi.PingStub, error) {
	args := this.Args()
	addr := "gopi"
	if len(args) == 1 {
		addr = args[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if conn, err := this.ConnPool.ConnectService(ctx, "tcp", addr, 0); err != nil {
		return nil, err
	} else if stub, _ := conn.NewStub("gopi.ping.Ping").(gopi.PingStub); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub")
	} else {
		return stub, nil
	}
}

func (this *app) GetMetricsStub() (gopi.MetricsStub, error) {
	args := this.Args()
	addr := "grpc"
	if len(args) == 1 {
		addr = args[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if conn, err := this.ConnPool.ConnectService(ctx, "tcp", addr, 0); err != nil {
		return nil, err
	} else if stub, _ := conn.NewStub("gopi.metrics.Metrics").(gopi.MetricsStub); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub")
	} else {
		return stub, nil
	}
}

func (this *app) GetServeAddress() (string, string, error) {
	var network, addr string

	args := this.Args()
	switch {
	case len(args) == 0:
		network = "tcp"
		addr = ":0"
	case len(args) == 1:
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			return "", "", err
		} else {
			network = "tcp"
			addr = args[0]
		}
	default:
		return "", "", gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	}

	// Return success
	return network, addr, nil
}
