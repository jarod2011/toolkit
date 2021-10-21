package buffer

import (
	"errors"
)

// ErrBufferCapacityNotEnough is defined error of buffer full
// When write bytes to buffer and buffer remain space not enough this error will throw.
var ErrBufferCapacityNotEnough = errors.New("buffer capacity is not enough")

// Buffer is bytes buffer interface
type Buffer interface {
	// NextN is get next n bytes in buffer
	// This method will return count and next bytes.
	// This method not move bytes pointer.
	NextN(n int) (int, []byte)
	// ShiftN is move bytes pointer forward
	ShiftN(n int) int
	// ReadN is get next bytes in buffer and move bytes pointer forward.
	ReadN(n int) (int, []byte)
	// Size will return bytes length in buffer
	Size() int
	// Capacity will return the buffer capacity
	Capacity() int
	// Write will write bytes to buffer
	// When buffer full will throw ErrBufferCapacityNotEnough error
	Write(p []byte) (int, error)
	// Reset will reset the buffer
	Reset()
	// Bytes will return all bytes in buffer
	Bytes() []byte
}

// DefaultBufferCapacity is default buffer Capacity
const DefaultBufferCapacity = 4096

type ringBuffer struct {
	start    int
	end      int
	capacity int
	buf      []byte
	full     bool
}

func (r *ringBuffer) NextN(n int) (int, []byte) {
	if n <= 0 {
		n = 0
	}
	if n == 0 {
		return 0, []byte{}
	}
	size := r.Size()
	if size == 0 {
		return 0, []byte{}
	}
	if n > size {
		n = size
	}
	if m := r.start + n; m < r.capacity {
		return n, r.buf[r.start:m][:]
	}
	return n, append(append([]byte{}, r.buf[r.start:]...), r.buf[0:n-(r.capacity-r.start)]...)
}

func (r *ringBuffer) ShiftN(n int) int {
	if n <= 0 {
		return 0
	}
	size := r.Size()
	if size < n {
		n = size
	}
	r.start += n
	if r.start >= r.capacity {
		r.start -= r.capacity
	}
	if r.full {
		r.full = false
	}
	return n
}

func (r *ringBuffer) ReadN(n int) (int, []byte) {
	n, b := r.NextN(n)
	r.ShiftN(n)
	return n, b
}

func (r *ringBuffer) Size() int {
	if r.start == r.end {
		if r.full {
			return r.capacity
		}
		return 0
	} else if r.start < r.end {
		return r.end - r.start
	}
	return r.capacity - r.start + r.end
}

func (r *ringBuffer) Capacity() int {
	return r.capacity
}

func (r *ringBuffer) Write(p []byte) (int, error) {
	pl := len(p)
	var err error
	rl := r.capacity - r.Size()
	if rl < pl {
		err = ErrBufferCapacityNotEnough
		pl = rl
	}
	if pl == 0 {
		return 0, err
	}
	if r.start > r.end {
		copy(r.buf[r.end:r.start], p[0:pl])
	} else {
		if m := r.capacity - r.end; m < pl {
			copy(r.buf[r.end:], p[0:m])
			copy(r.buf[0:r.start], p[m:])
		} else {
			copy(r.buf[r.end:], p[0:pl])
		}
	}
	r.end += pl
	if r.end >= r.capacity {
		r.end -= r.capacity
	}
	if pl == rl {
		r.full = true
	}
	return pl, err
}

func (r *ringBuffer) Reset() {
	r.start = 0
	r.end = 0
}

func (r *ringBuffer) Bytes() []byte {
	_, b := r.NextN(r.Size())
	return b
}

// NewBuffer will create buffer of capacity
func NewBuffer(capacity int) *ringBuffer {
	if capacity <= 0 {
		capacity = DefaultBufferCapacity
	}
	return &ringBuffer{
		start:    0,
		end:      0,
		capacity: capacity,
		buf:      make([]byte, capacity),
	}
}
