package server

import "github.com/jarod2011/toolkit/buffer"

// Codec will Decode receive data from client and Encode send data before write to client
type Codec interface {
	// Encode will encode data before write to client
	Encode([]byte) []byte
	// Decode will decode receive data from client
	Decode(buf buffer.Buffer) []byte
}
