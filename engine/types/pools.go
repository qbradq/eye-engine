package types

import (
	"errors"
	"sync"
)

// pointI2DPool is the pool of PointI2D slices
var pointI2DPool sync.Pool = sync.Pool{
	New: func() any {
		b := make([]PointI2D, 0, 512)
		return &b
	},
}

// GetPointI2dPoolSlice returns a possibly reused, empty PointI2D slice.
func GetPointI2dPoolSlice() []PointI2D {
	retPtr, ok := pointI2DPool.Get().(*[]PointI2D)
	if !ok {
		panic(errors.New("PointI2DPool.Get() did not return *[]PointI2D"))
	}
	ret := *retPtr
	return ret
}

// ReleasePointI2DPoolSlice returns s to the available memory pool.
func ReleasePointI2DPoolSlice(s []PointI2D) {
	sPtr := &s
	*sPtr = (*sPtr)[:0]
	pointI2DPool.Put(sPtr)
}
