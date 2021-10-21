package server

import "github.com/jarod2011/toolkit/logger"

// Options defined server options
type Options struct {
	// Logger is Logger implements
	Logger logger.Logger
	// Task is Task implements function
	// This function will call when server started
	Task Task
	// Codec is Codec implements
	// All data receive and send will use Codec Encode and Decode
	Codec Codec
}

type Option func(options *Options)

// WithLogger is edit Options Logger field
func WithLogger(logger logger.Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

// WithTask is edit Options Task field
func WithTask(task Task) Option {
	return func(options *Options) {
		options.Task = task
	}
}

// WithCodec is edit Options Codec field
func WithCodec(codec Codec) Option {
	return func(options *Options) {
		options.Codec = codec
	}
}
