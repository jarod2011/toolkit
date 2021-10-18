package mq

// Options is used to extend for subscribe method
type Options map[string]interface{}

// Option is callback function to edit Options
type Option func(opts *Options)
