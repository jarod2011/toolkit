package server

import "github.com/jarod2011/toolkit/logger"

type Connection interface {
	Send(data []byte, withoutEncode bool) error
	Remote() string
	Local() string
	Logger() logger.Logger
}
