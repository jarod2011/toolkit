package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jarod2011/toolkit/logger"
)

func TestWithLogger(t *testing.T) {
	options := Options{}
	assert.Nil(t, options.Logger)
	WithLogger(logger.NewLogger())(&options)
	assert.NotNil(t, options.Logger)
}

func TestWithTask(t *testing.T) {
	options := Options{}
	assert.Nil(t, options.Task)
	WithTask(func() (time.Duration, Action) {
		return 0, 0
	})(&options)
	assert.NotNil(t, options.Task)
}

func TestWithCodec(t *testing.T) {
	options := Options{}
	assert.Nil(t, options.Codec)
	WithCodec(new(NothingCodec))(&options)
	assert.NotNil(t, options.Codec)
}
