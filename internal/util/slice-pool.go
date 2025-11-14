package util

import (
	"errors"
	"sync"
)

// SlicePool manages a pool of typed slices.
type SlicePool[T any] struct {
	p sync.Pool
}

// NewSlicePool returns a new SlicePool that generates new slices with capacity
// set to cap.
func NewSlicePool[T any](cap int) *SlicePool[T] {
	return &SlicePool[T]{
		p: sync.Pool{
			New: func() any {
				b := make([]T, 0, cap)
				return &b
			},
		},
	}
}

// Get returns a possibly reused, empty typed slice.
func (p *SlicePool[T]) Get() []T {
	retPtr, ok := p.p.Get().(*[]T)
	if !ok {
		panic(errors.New("SlicePool.Get() did not return *[]T"))
	}
	ret := *retPtr
	return ret
}

// Release returns s to the available memory pool.
func (p *SlicePool[T]) Release(s []T) {
	sPtr := &s
	*sPtr = (*sPtr)[:0]
	p.p.Put(sPtr)
}
