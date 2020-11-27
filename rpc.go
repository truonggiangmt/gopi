<<<<<<< HEAD
=======
/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

>>>>>>> master
package gopi

import "context"

/////////////////////////////////////////////////////////////////////
// INTERFACES

// Server is a generic gRPC server, which can serve registered services
type Server interface {
	RegisterService(interface{}, Service) error   // Register an RPC service
	StartInBackground(network, addr string) error // Start server in background and return
	Stop(bool) error                              // Stop server, when argument is true forcefully disconnects any clients
	Addr() string                                 // Addr returns the address of the server, or empty if not connected
}

// Service defines an RPC service, which can cancel any on-going streams
// when server stops
type Service interface {
	CancelStreams()
}

// ConnPool is a factory of client connections
type ConnPool interface {
	Connect(network, addr string) (Conn, error)
}

// Conn is a connection to a remote server
type Conn interface {
	// Addr returns the bound server address, or empty string if connection is closed
	Addr() string

	// Mutex
	Lock()   // Lock during RPC call
	Unlock() // Unlock at end of RPC call

	// Methods
	ListServices(context.Context) ([]string, error) // Return a list of services supported
	NewStub(string) ServiceStub                     // Return the stub for a named service
}

// ServiceStub is a client-side stub used to invoke remote service methods
type ServiceStub interface {
	New(Conn)
}

/////////////////////////////////////////////////////////////////////
// SERVICES

type PingService interface {
	Service
}

type PingStub interface {
	ServiceStub

	Ping(ctx context.Context) error
	Version(ctx context.Context) (Version, error)
}
