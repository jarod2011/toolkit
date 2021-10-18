package mq

import "time"

// Options is used to extend for subscribe method
type Options map[string]interface{}

// Option is callback function to edit Options
type Option func(opts *Options)

const (
	TopicOptionKey   = "topic"   // the topic key name in options
	TimeoutOptionKey = "timeout" // the timeout key name in options
)

// WithTopic is set message topic or set subscriber only subscribe some topic
func WithTopic(topic string) Option {
	return func(opts *Options) {
		(*opts)[TopicOptionKey] = topic
	}
}

// WithTimeout is set subscriber receive message timeout
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		(*opts)[TimeoutOptionKey] = timeout
	}
}
