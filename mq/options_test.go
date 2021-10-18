package mq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithTopic(t *testing.T) {
	var opts Options = make(map[string]interface{})
	_, ok := opts[TopicOptionKey]
	assert.False(t, ok)
	WithTopic("hello")(&opts)
	topic, ok := opts[TopicOptionKey]
	assert.Equal(t, topic, "hello")
	assert.True(t, ok)
}

func TestWithTimeout(t *testing.T) {
	var opts Options = make(map[string]interface{})
	_, ok := opts[TimeoutOptionKey]
	assert.False(t, ok)
	WithTimeout(time.Minute)(&opts)
	timeout, ok := opts[TimeoutOptionKey]
	assert.Equal(t, timeout.(time.Duration), time.Minute)
	assert.True(t, ok)
}
