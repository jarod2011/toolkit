package buffer

import "sync"

// Pool is buffer pool
type Pool struct {
	po *sync.Pool
}

// Get is get Buffer from pool
func (p *Pool) Get() Buffer {
	return p.po.Get().(Buffer)
}

// Put is put Buffer to pool
func (p *Pool) Put(buf Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	p.po.Put(buf)
}

// NewPool will get a Buffer pool
func NewPool(bufferCapacity int) *Pool {
	return &Pool{
		po: &sync.Pool{
			New: func() interface{} {
				return NewBuffer(bufferCapacity)
			},
		},
	}
}
