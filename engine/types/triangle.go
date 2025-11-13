package types

import (
	"sort"
)

// TriangleScanLines returns a slice of an even number of points. For each pair
// of points, they specify the left and right limits of a scan line contained
// within the triangle.
func TriangleScanLines(p0, p1, p2 PointI2D) []PointI2D {
	ret := GetPointI2dPoolSlice()
	// Sort points top to bottom
	p := []PointI2D{p0, p1, p2}
	sort.SliceStable(p, func(i, j int) bool {
		return p[i][1] < p[j][1]
	})
	p0 = p[0]
	p1 = p[1]
	p2 = p[2]
	// Setup the interpolation loops
	s0 := float32(p2[0]-p1[0]) / float32(p2[1]-p1[1])
	s1 := float32(p1[0]-p0[0]) / float32(p1[1]-p0[1])
	s2 := float32(p2[0]-p0[0]) / float32(p2[1]-p0[1])
	xLeft := float32(p0[0])
	xRight := float32(p0[0])
	// Top interpolation loop
	lSlope := s1
	rSlope := s2
	if s1 > s2 {
		lSlope = s2
		rSlope = s1
	}
	for y := p0[1]; y < p1[1]; y++ {
		ret = append(ret, PointI2D{int(xLeft), y}, PointI2D{int(xRight), y})
		xLeft += lSlope
		xRight += rSlope
	}
	// Bottom interpolation loop
	if lSlope == s2 {
		rSlope = s0
	} else {
		lSlope = s0
	}
	for y := p1[1]; y <= p2[1]; y++ {
		ret = append(ret, PointI2D{int(xLeft), y}, PointI2D{int(xRight), y})
		xLeft += lSlope
		xRight += rSlope
	}
	return ret
}
