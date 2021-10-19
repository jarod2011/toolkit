package mq

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryBroker(t *testing.T) {
	mem := NewMemoryBroker(context.TODO(), 1)
	assert.Implements(t, (*Broker)(nil), mem)
}

func TestMemoryBroker_Close(t *testing.T) {
	t.Run("test close by context cancel", func(t *testing.T) {
		r1 := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		mem := NewMemoryBroker(ctx, 5)
		r2 := runtime.NumGoroutine()
		assert.Less(t, r1, r2)
		cancel()
		time.Sleep(time.Millisecond)
		assert.True(t, mem.closed == 1)
		r3 := runtime.NumGoroutine()
		assert.True(t, r1 == r3)
	})
	t.Run("test close method", func(t *testing.T) {
		r1 := runtime.NumGoroutine()
		mem := NewMemoryBroker(context.TODO(), 5)
		r2 := runtime.NumGoroutine()
		assert.Less(t, r1, r2)
		mem.Close()
		time.Sleep(time.Millisecond)
		assert.True(t, mem.closed == 1)
		r3 := runtime.NumGoroutine()
		assert.True(t, r3 < r2)
	})
	t.Run("test GC", func(t *testing.T) {
		r1 := runtime.NumGoroutine()
		mem := NewMemoryBroker(context.TODO(), 5)
		r2 := runtime.NumGoroutine()
		assert.Less(t, r1, r2)
		_ = mem
		mem = nil
		runtime.GC()
		time.Sleep(time.Millisecond)
		r3 := runtime.NumGoroutine()
		assert.True(t, r3 < r2)
	})
}

func TestMemoryBroker_Publish(t *testing.T) {
	mem := NewMemoryBroker(context.TODO(), 0)
	defer func() {
		assert.ErrorIs(t, mem.Publish(555), ErrBrokerClosed)
	}()
	defer mem.Close()
	sub1 := make(chan interface{}, 10)
	defer close(sub1)
	mem.Subscribe(sub1, WithTopic("event"))
	defer mem.Unsubscribe(sub1)
	sub2 := make(chan interface{}, 10)
	defer close(sub2)
	mem.Subscribe(sub2)
	defer mem.Unsubscribe(sub2)
	t.Run("test publish without topic", func(t *testing.T) {
		mem.Publish(123)
		s2 := <-sub2
		assert.Equal(t, s2, 123)
		select {
		case s1 := <-sub1:
			t.Errorf("subscribe 1 only subscribe topic 'event', so should not get this message: %v", s1)
		default:
		}
	})
	t.Run("test publish with topic", func(t *testing.T) {
		mem.Publish(234, WithTopic("event"))
		s1 := <-sub1
		assert.Equal(t, s1, 234)
		s2 := <-sub2
		assert.Equal(t, s2, 234)
		mem.Publish(345, WithTopic("event1"))
		s3 := <-sub2
		assert.Equal(t, s3, 345)
		select {
		case s1 := <-sub1:
			t.Errorf("subscribe 1 only subscribe topic 'event', so should not get this message: %v", s1)
		default:
		}
	})
	t.Run("test publish with timeout", func(t *testing.T) {
		sub3 := make(chan interface{})
		defer close(sub3)
		mem.Subscribe(sub3, WithTimeout(time.Millisecond*500))
		defer mem.Subscribe(sub3)
		mem.Publish(456)
		go func() {
			<-sub1
		}()
		go func() {
			<-sub2
		}()
		time.Sleep(time.Second)
		select {
		case s3 := <-sub3:
			t.Errorf("subscribe 3 should get message timeout, but get message: %v", s3)
		default:
		}
	})
}

func TestMemoryBroker_Subscribe(t *testing.T) {
	mem := NewMemoryBroker(context.TODO(), 0)
	assert.Nil(t, mem.Subscribe(make(chan interface{})))
	mem.Close()
	assert.ErrorIs(t, mem.Subscribe(make(chan interface{})), ErrBrokerClosed)
}

func TestMemoryBroker_Unsubscribe(t *testing.T) {
	mem := NewMemoryBroker(context.TODO(), 0)
	assert.Nil(t, mem.Unsubscribe(make(chan interface{})))
	mem.Close()
	assert.ErrorIs(t, mem.Unsubscribe(make(chan interface{})), ErrBrokerClosed)
}
