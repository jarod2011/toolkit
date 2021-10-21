package server

// Handler is handler of client read or write data
// The Action of method returned can control server or connection
// When Action is StopServerAction: the server will stop
// When Action is DisconnectionAction: the connection will be close
type Handler interface {

	// OnConnected will call when client connected
	// the result action can control server
	OnConnected(conn Connection) (Action, error)

	// OnDisconnected will call when client disconnected
	OnDisconnected(conn Connection) error

	// OnReceived will call when read data from client
	// If Codec is config, decode data(frame) will be input when frame not nil
	// the result action can control server
	OnReceived(frame []byte, conn Connection) (Action, error)

	// OnError will call when any error occurred
	// the result action can control server
	OnError(conn Connection, err error) Action
}
