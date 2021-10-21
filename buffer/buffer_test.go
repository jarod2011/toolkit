package buffer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer(10)
	assert.Equal(t, buf.capacity, 10)
	buf = NewBuffer(-1)
	assert.Equal(t, buf.capacity, DefaultBufferCapacity)
}

func TestRingBuffer_Capacity(t *testing.T) {
	buf := NewBuffer(100)
	assert.Equal(t, buf.Capacity(), 100)
	buf = NewBuffer(0)
	assert.Equal(t, buf.Capacity(), DefaultBufferCapacity)
}

func TestRingBuffer_Size(t *testing.T) {
	buf := NewBuffer(10)
	assert.Equal(t, buf.Size(), 0)
	buf.Write([]byte{0x01, 0x02})
	assert.Equal(t, buf.Size(), 2)
	buf.Write([]byte{0x03, 0x04, 0x05})
	assert.Equal(t, buf.Size(), 5)
	buf.ShiftN(2)
	assert.Equal(t, buf.Size(), 3)
	buf.Write([]byte{0x06, 0x07, 0x08, 0x09, 0x0a})
	assert.Equal(t, buf.Size(), 8)
	assert.Equal(t, buf.start, 2)
	assert.Equal(t, buf.end, 0)
	buf.Write([]byte{0x0b})
	assert.Equal(t, buf.Size(), 9)
	assert.Equal(t, buf.end, 1)
	assert.Equal(t, buf.full, false)
	buf.Write([]byte{0x0c})
	assert.Equal(t, buf.Size(), 10)
	assert.Equal(t, buf.end, 2)
	assert.Equal(t, buf.full, true)
}

func TestRingBuffer_NextN(t *testing.T) {
	buf := NewBuffer(10)
	n, b := buf.NextN(-1)
	assert.Equal(t, n, 0)
	assert.Empty(t, b)
	n, b = buf.NextN(10)
	assert.Equal(t, n, 0)
	assert.Empty(t, b)
	n, err := buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05})
	assert.Nil(t, err)
	assert.Equal(t, n, 5)
	n, b = buf.NextN(6)
	assert.Equal(t, n, 5)
	assert.EqualValues(t, b, []byte{0x01, 0x02, 0x03, 0x04, 0x05})
	n, err = buf.Write([]byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b})
	assert.Equal(t, n, 5)
	assert.ErrorIs(t, err, ErrBufferCapacityNotEnough)
	n, b = buf.NextN(11)
	assert.Equal(t, n, 10)
	assert.EqualValues(t, b, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a})
	buf.ShiftN(5)
	n, err = buf.Write([]byte{0x0b, 0x0c, 0x0d, 0x0e})
	assert.Equal(t, n, 4)
	assert.Nil(t, err)
	n, b = buf.NextN(8)
	assert.Equal(t, n, 8)
	assert.EqualValues(t, b, []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d})
	n, err = buf.Write([]byte{})
	assert.Equal(t, n, 0)
	assert.Nil(t, err)
	buf.ShiftN(10)
	assert.Equal(t, buf.Size(), 0)
	n, err = buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	assert.Equal(t, n, 8)
	assert.Nil(t, err)
	n, b = buf.NextN(10)
	assert.Equal(t, n, 8)
	assert.EqualValues(t, b, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
}

func TestRingBuffer_ShiftN(t *testing.T) {
	buf := NewBuffer(10)
	assert.Equal(t, buf.Size(), 0)
	buf.Write([]byte{0x01, 0x02})
	assert.Equal(t, buf.Size(), 2)
	buf.ShiftN(1)
	assert.Equal(t, buf.Size(), 1)
	n, b := buf.NextN(2)
	assert.Equal(t, n, 1)
	assert.EqualValues(t, b, []byte{0x02})
	n = buf.ShiftN(-1)
	assert.Equal(t, n, 0)
	n = buf.ShiftN(10)
	assert.Equal(t, n, 1)
	buf.start = 8
	buf.end = 4
	n = buf.ShiftN(4)
	assert.Equal(t, n, 4)
	buf.Write(bytes.Repeat([]byte{0x01}, 10))
	assert.Equal(t, buf.full, true)
	buf.ShiftN(1)
	assert.Equal(t, buf.full, false)
}

func TestRingBuffer_ReadN(t *testing.T) {
	buf := NewBuffer(10)
	buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	n, b := buf.ReadN(3)
	assert.Equal(t, n, 3)
	assert.EqualValues(t, b, []byte{0x01, 0x02, 0x03})
	n, b = buf.ReadN(10)
	assert.Equal(t, n, 3)
	assert.EqualValues(t, b, []byte{0x04, 0x05, 0x06})
}

func TestRingBuffer_Bytes(t *testing.T) {
	buf := NewBuffer(10)
	buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	assert.EqualValues(t, buf.Bytes(), []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	buf.Write([]byte{0x07, 0x08, 0x09, 0x0a, 0x0b})
	assert.EqualValues(t, buf.Bytes(), []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a})
}

func TestRingBuffer_Reset(t *testing.T) {
	buf := NewBuffer(10)
	buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	assert.EqualValues(t, buf.Bytes(), []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	buf.Reset()
	assert.EqualValues(t, buf.Bytes(), []byte{})
}

func TestRingBuffer_Write(t *testing.T) {
	buf := NewBuffer(10)
	n, err := buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	assert.Equal(t, n, 6)
	assert.Nil(t, err)
	n, err = buf.Write([]byte{0x07, 0x08, 0x09, 0x0a, 0x0b})
	assert.Equal(t, n, 4)
	assert.ErrorIs(t, err, ErrBufferCapacityNotEnough)
	n, b := buf.ReadN(5)
	assert.Equal(t, n, 5)
	assert.EqualValues(t, b, []byte{0x01, 0x02, 0x03, 0x04, 0x05})
	n, err = buf.Write([]byte{0x0b, 0x0c, 0x0d})
	assert.Equal(t, n, 3)
	assert.EqualValues(t, buf.Bytes(), []byte{0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d})
	assert.Equal(t, buf.Size(), 8)
	n, err = buf.Write([]byte{0x0e, 0x0f, 0x10, 0x11})
	assert.Equal(t, n, 2)
	assert.ErrorIs(t, err, ErrBufferCapacityNotEnough)
	n, err = buf.Write([]byte{0x12, 0x13})
	assert.Equal(t, n, 0)
	assert.ErrorIs(t, err, ErrBufferCapacityNotEnough)
	buf.start = 6
	buf.end = 8
	n, err = buf.Write([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	assert.Equal(t, n, 6)
	assert.Nil(t, err)
}
