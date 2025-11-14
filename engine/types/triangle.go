package types

import (
	"fmt"

	"github.com/qbradq/eye-engine/internal/util"
)

// ClipPolygonToPlane appends the points of the polygon in to out clipped to the
// named plane of the given rect.
func ClipPolygonToPlane(in, out []PointI2D, r RectI2D, oc OutCode) {
	s := in[len(in)-1]
	for _, p := range in {
		sIn := ((r.OutCode(s)) & oc) == 0
		pIn := ((r.OutCode(p)) & oc) == 0
		dy := p[1] - s[1]
		dx := p[0] - s[0]
		var pi PointI2D
		switch oc {
		case OutCodeLeft:
			xb := r[0]
			ratio := float32(xb-s[0]) / float32(dx)
			yi := float32(s[1]) + float32(dy)*ratio
			pi = PointI2D{xb, int(yi)}
		case OutCodeRight:
			xb := r[0] + r[2] - 1
			ratio := float32(xb-s[0]) / float32(dx)
			yi := float32(s[1]) + float32(dy)*ratio
			pi = PointI2D{xb, int(yi)}
		case OutCodeTop:
			yb := r[1]
			ratio := float32(yb-s[1]) / float32(dy)
			xi := float32(s[0]) + float32(dx)*ratio
			pi = PointI2D{int(xi), yb}
		case OutCodeBottom:
			yb := r[1] + r[3] - 1
			ratio := float32(yb-s[1]) / float32(dy)
			xi := float32(s[0]) + float32(dx)*ratio
			pi = PointI2D{int(xi), yb}
		default:
			panic(fmt.Errorf("unhandled out code %d", oc))
		}
		if sIn && pIn {
			out = append(out, p)
		} else if sIn && !pIn {
			out = append(out, pi)
		} else if !sIn && pIn {
			out = append(out, pi, p)
		} // else trivial reject, output nothing
		s = p
	}
}

// ClipPolygonToRect appends the points of polygon clipped to r to out.
func ClipPolygonToRect(polygon, out []PointI2D, r RectI2D) {
	b1 := pointI2DPool.Get()
	b2 := pointI2DPool.Get()
	ClipPolygonToPlane(polygon, b1, r, OutCodeLeft)
	ClipPolygonToPlane(b1, b2, r, OutCodeRight)
	b1 = b1[:0]
	ClipPolygonToPlane(b2, b1, r, OutCodeBottom)
	ClipPolygonToPlane(b1, out, r, OutCodeTop)
	pointI2DPool.Release(b1)
	pointI2DPool.Release(b2)
}

// polyFillEdgeDef defines an edge used during polygon filling.
type polyFillEdgeDef struct {
	y int     // Starting Y position of the edge
	x float32 // Current X position of the edge
	m float32 // Inverse slope
}

// PolygonScanLines returns a slice of an even number of points. For each pair
// of points, they specify the left and right limits of a scan line contained
// within polygon clipped to r. When done with the returned slice, pass it to
// ReleasePointI2DPool() to release it.
func PolygonScanLines(polygon []PointI2D, r RectI2D) []PointI2D {
	if len(polygon) < 3 {
		return nil
	}
	ret := pointI2DPool.Get()
	// Find min and max Y
	minY := polygon[0][1]
	maxY := polygon[0][1]
	for _, p := range polygon {
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	// Build edge table
	edges := make([][]*polyFillEdgeDef, (maxY - minY))
	pa := polygon[len(polygon)-1]
	for _, pb := range polygon {
		min := pa
		max := pb
		if min[1] > max[1] {
			t := min
			min = max
			max = t
		}
		dx := max[0] - min[0]
		dy := max[1] - min[1]
		if dy < 1 {
			pa = pb
			continue
		}
		edge := &polyFillEdgeDef{
			y: min[1],
			x: float32(min[0]),
			m: float32(dx) / float32(dy),
		}
		for iy := min[1]; iy < max[1]; iy++ {
			edges[iy-minY] = append(edges[iy-minY], edge)
		}
		pa = pb
	}
	// Scan line loop
	var ael []*polyFillEdgeDef
	for y := minY; y < maxY; y++ {
		idx := y - minY
		ael = edges[idx]
		l := ael[0]
		r := ael[0]
		for _, v := range ael {
			if v.x < l.x {
				l = v
			}
			if v.x > r.x {
				r = v
			}
		}
		ret = append(
			ret,
			PointI2D{int(l.x), l.y},
			PointI2D{int(r.x), r.y},
		)
		ael[0].x += ael[0].m
		ael[1].x += ael[1].m
		ael[0].y++
		ael[1].y++
	}
	return ret
}

// TriangleScanLines returns a slice of an even number of points. For each pair
// of points, they specify the left and right limits of a scan line contained
// within the triangle. All scan lines are clipped to r. When done with the
// returned slice, pass it to ReleasePointI2DPoolSlice() to release it.
func TriangleScanLines(p []PointI2D, r RectI2D) []PointI2D {
	ret := pointI2DPool.Get()
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
