package mq

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

// MemoryBroker is a Broker implement by Memory use map and channel
type MemoryBroker struct {
	*memoryBroker
}

type subscribeItem struct {
	ch   chan<- interface{}
	opts Options
}

type publishItem struct {
	value interface{}
	opts  Options
}

type memoryBroker struct {
	subscribeChannel   chan subscribeItem
	unsubscribeChannel chan chan<- interface{}
	publishChannel     chan publishItem
	subscribers        map[chan<- interface{}]Options
	closed             uint32
	ctx                context.Context
	cancel             context.CancelFunc
}

// NewMemoryBroker is build a Broker implement by memory
func NewMemoryBroker(ctx context.Context, capacity int) *MemoryBroker {
	m := &MemoryBroker{
		&memoryBroker{
			subscribeChannel:   make(chan subscribeItem),
			unsubscribeChannel: make(chan chan<- interface{}),
			publishChannel:     make(chan publishItem, capacity),
			subscribers:        make(map[chan<- interface{}]Options),
			closed:             0,
		},
	}
	m.memoryBroker.ctx, m.memoryBroker.cancel = context.WithCancel(ctx)
	go m.loop()
	runtime.SetFinalizer(m, func(obj *MemoryBroker) {
		obj.Close()
	})
	return m
}

func (m *memoryBroker) opened() bool {
	return atomic.LoadUint32(&m.closed) == 0
}

func (m *memoryBroker) loop() {
	defer func() {
		close(m.subscribeChannel)
		close(m.unsubscribeChannel)
	}()
	for {
		select {
		case <-m.ctx.Done():
			m.Close()
		case sub := <-m.subscribeChannel:
			m.subscribers[sub.ch] = sub.opts
		case uns := <-m.unsubscribeChannel:
			delete(m.subscribers, uns)
		case item, ok := <-m.publishChannel:
			if !ok {
				return
			}
			m.publishToAllSubscribers(item)
		}
	}
}

func (m *memoryBroker) publishToAllSubscribers(item publishItem) {
	for ch, opts := range m.subscribers {
		if onlyTopic, ok := opts[TopicOptionKey]; ok {
			if topic, ok := item.opts[TopicOptionKey]; ok {
				if onlyTopic.(string) == topic.(string) {
					goto pub
				}
			}
			continue
		}
	pub:
		ctx, cancel := context.WithCancel(m.ctx)
		if timeout, ok := opts[TimeoutOptionKey]; ok {
			ctx, cancel = context.WithTimeout(m.ctx, timeout.(time.Duration))
		}
		select {
		case ch <- item.value:
		case <-ctx.Done():
		}
		cancel()
	}
}

// Publish is publish to memory broker
func (m *memoryBroker) Publish(v interface{}, opts ...Option) error {
	if !m.opened() {
		return ErrBrokerClosed
	}
	item := publishItem{
		value: v,
		opts:  make(map[string]interface{}),
	}
	for _, o := range opts {
		o(&item.opts)
	}
	m.publishChannel <- item
	return nil
}

// Subscribe is subscribe from memory broker
func (m *memoryBroker) Subscribe(ch chan<- interface{}, opts ...Option) error {
	if !m.opened() {
		return ErrBrokerClosed
	}
	item := subscribeItem{
		ch:   ch,
		opts: make(map[string]interface{}),
	}
	for _, o := range opts {
		o(&item.opts)
	}
	m.subscribeChannel <- item
	return nil
}

// Unsubscribe is unsubscribe from memory broker
func (m *memoryBroker) Unsubscribe(ch chan<- interface{}) error {
	if !m.opened() {
		return ErrBrokerClosed
	}
	m.unsubscribeChannel <- ch
	return nil
}

// Close is stop the memory broker
func (m *memoryBroker) Close() error {
	if !atomic.CompareAndSwapUint32(&m.closed, 0, 1) {
		return nil
	}
	close(m.publishChannel)
	m.cancel()
	return nil
}
