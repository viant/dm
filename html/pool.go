package html

import (
	"sync"
	"sync/atomic"
)

// Pool represents sharable DOM instances
type Pool struct {
	pool    *sync.Pool
	counter int32
	maxSize int32
	dom     *VirtualDOM
}

// NewPool creates DOM pool
func NewPool(size int32, dom *VirtualDOM) *Pool {
	return &Pool{
		pool: &sync.Pool{
			New: func() interface{} {
				return NewBuffer(dom.initialBufferSize)
			},
		},
		dom:     dom,
		counter: 0,
		maxSize: size,
	}
}

// Put returns DOM to the pool
func (p *Pool) Put(dom *DOM) {
	diff := int32(atomic.AddInt32(&p.counter, -1))
	if diff < p.maxSize {
		p.pool.Put(dom.buffer)
	}
}

// New creates or reuse recent DOM instance
func (p *Pool) New() *DOM {
	atomic.AddInt32(&p.counter, 1)
	return p.dom.DOM(p.pool.New().(*Buffer))
}
