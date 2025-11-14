package types

import (
	"math"

	"github.com/qbradq/eye-engine/internal/util"
)

// TriangleScanLines returns a slice of an even number of points. For each pair
// of points, they specify the left and right limits of a scan line contained
// within the triangle. All scan lines are clipped to r. When done with the
// returned slice, pass it to ReleasePointI2DPoolSlice() to release it.
func TriangleScanLines(p []PointI2D, r RectI2D) []PointI2D {
	ret := PointI2DPool.Get()
	rl := r[0]
	rt := r[1]
	rr := r[0] + r[2] - 1
	rb := r[1] + r[3] - 1
	// Unrolled insertion sort, n=3
	if p[0][1] > p[1][1] {
		p[0], p[1] = p[1], p[0]
	}
	if p[1][1] > p[2][1] {
		p[1], p[2] = p[2], p[1]
		if p[0][1] > p[1][1] {
			p[0], p[1] = p[1], p[0]
		}
	}
	// Calculate inverse slopes
	s0 := float32(p[2][0]-p[1][0]) / float32(p[2][1]-p[1][1])
	s1 := float32(p[1][0]-p[0][0]) / float32(p[1][1]-p[0][1])
	s2 := float32(p[2][0]-p[0][0]) / float32(p[2][1]-p[0][1])
	xLeft := float32(p[0][0])
	xRight := float32(p[0][0])
	// Interpolation loop setup
	lSlope := s1
	rSlope := s2
	if s1 > s2 {
		lSlope = s2
		rSlope = s1
	}
	// Clip vertical
	tt := p[0][1]
	tb := p[2][1]
	yMin := util.MaxInt(rt, tt)
	yMax := util.MinInt(rb, tb)
	if yMin > tt {
		xLeft += lSlope * float32(yMin-tt)
		xRight += rSlope * float32(yMin-tt)
	}
	// Trivial discard case
	if tb < yMin || tt > yMax {
		return nil
	}
	// Top interpolation loop
	yTopLoopEnd := util.MinInt(p[1][1], yMax)
	for y := yMin; y < yTopLoopEnd; y++ {
		// Trivial discard case
		l := int(xLeft)
		r := int(xRight)
		xLeft += lSlope
		xRight += rSlope
		if r < rl || l > rr {
			continue
		}
		// Clip horizontal
		ret = append(
			ret,
			PointI2D{util.MaxInt(l, rl), y},
			PointI2D{util.MinInt(r, rr), y},
		)
	}
	// Horizontal line hackery
	if math.IsInf(float64(lSlope), 0) {
		xLeft = float32(p[1][0])
	}
	if math.IsInf(float64(rSlope), 0) {
		xRight = float32(p[1][0])
	}
	// Bottom interpolation loop
	if lSlope == s2 {
		rSlope = s0
	} else {
		lSlope = s0
	}
	yBottomStart := util.MaxInt(p[1][1], yMin)
	for y := yBottomStart; y <= yMax; y++ {
		// Trivial discard case
		l := int(xLeft)
		r := int(xRight)
		if r < rl || l > rr {
			continue
		}
		xLeft += lSlope
		xRight += rSlope
		// Clip horizontal
		ret = append(
			ret,
			PointI2D{util.MaxInt(l, rl), y},
			PointI2D{util.MinInt(r, rr), y},
		)
	}
	return ret
}
