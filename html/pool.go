package html

import (
	"sync"
	"sync/atomic"
)

//Pool represents sharable Document instances
type Pool struct {
	pool    *sync.Pool
	counter int32
	maxSize int
	dom     *DOM
}

//NewPool creates Document pool
func NewPool(size int, dom *DOM) *Pool {
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

//Put returns Document to the pool
func (p *Pool) Put(dom *Document) {
	diff := int(atomic.AddInt32(&p.counter, -1))
	if diff < p.maxSize {
		p.pool.Put(dom.buffer)
	}
}

//New creates or reuse recent Document instance
func (p *Pool) New() *Document {
	atomic.AddInt32(&p.counter, 1)
	return p.dom.Document(p.pool.New().(*Buffer))
}
