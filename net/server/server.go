package server

import (
	"errors"
)

// Action is control server or connection next action
type Action int

const (
	NothingAction       Action = iota // nothing to do
	DisconnectionAction               // this action will close connection
	StopServerAction                  // this action will stop server
)

// ErrServerClosed will throw when server closed
var ErrServerClosed = errors.New("server closed")

type Address struct {
}

// Server is multi address handler server
type Server interface {
	// Bind is bind address and handler to server
	Bind(address *Address, handler Handler) error

	// Stop will stop input addresses
	// When addresses empty, the server will stop all bind address
	Stop(addresses ...*Address) error

	// Start is start the server and blocking.
	// When server all addresses stop, will throw ErrServerClosed
	Start() error
}
