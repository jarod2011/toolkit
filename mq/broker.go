package mq

import "errors"

var (
	// ErrBrokerClosed is defined broker close error.
	// After broker call Close method:
	// 1. Subscriber call Subscribe or Unsubscribe will get this error
	// 2. Publisher call Publish will get this error
	ErrBrokerClosed = errors.New("broker is closed")
)

// Broker is defined a message queue broker.
// The method Publish message to all Subscriber.
// All subscriber can use method Subscribe to subscribe message
// and use method Unsubscribe to unsubscribe message
// The method Close will close broker, all new publish will get ErrBrokerClosed error.
// When broker closed, all remaining message in broker will publish until broker empty.
type Broker interface {

	// Publish is publish message to all subscribers cross broker.
	// When broker closed will get ErrBrokerClosed error.
	Publish(v interface{}, opts ...Option) error

	// Subscribe is subscribe from broker
	// When broker closed will get ErrBrokerClosed error
	// The opts is used to extend for subscribe params, such as subscribe some topic
	Subscribe(ch chan<- interface{}, opts ...Option) error

	// Unsubscribe is unsubscribe from broker
	// When broker closed will get ErrBrokerClosed error
	// This method is idempotent.
	Unsubscribe(ch chan<- interface{}) error

	// Close the broker
	// This method is idempotent.
	// After call this method, Publish,Subscribe,Unsubscribe method will get ErrBrokerClosed error.
	Close() error
}
