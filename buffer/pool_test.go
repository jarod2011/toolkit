package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	po := NewPool(10)
	buf := po.Get()
	assert.Equal(t, buf.Capacity(), 10)
	po.Put(nil)
	po.Put(buf)
}
