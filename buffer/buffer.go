package buffer

import (
	"errors"
)

var ErrBufferCapacityNotEnough = errors.New("buffer capacity is not enough")

type Buffer interface {
	NextN(n int) (int, []byte)
	ShiftN(n int) int
	ReadN(n int) (int, []byte)
	Size() int
	Capacity() int
	Write(p []byte) (int, error)
	Reset()
	Bytes() []byte
}

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
