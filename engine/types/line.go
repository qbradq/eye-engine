package types

// Line returns a list of points as in Line(), bounded to the given rect.
// When done with the returned slice, release it with ReleasePointI2DPoolSlice.
func Line(p0, p1 PointI2D, r RectI2D) []PointI2D {
	// Setup
	oc0 := r.OutCode(p0)
	oc1 := r.OutCode(p1)
	if oc0|oc1 == 0 {
		// Trivial accept
		return UnclippedLine(p0, p1)
	}
	if oc0&oc1 != 0 {
		// Trivial reject
		return nil
	}
	// Clipping loop
	for {
		// Break cases
		var p *PointI2D
		var oc *OutCode
		if oc0|oc1 == 0 {
			// Trivial accept
			return UnclippedLine(p0, p1)
		}
		if oc0&oc1 != 0 {
			// Trivial reject
			return nil
		}
		// Clipping selection
		if oc0 != 0 {
			p = &p0
			oc = &oc0
		} else {
			p = &p1
			oc = &oc1
		}
		// Clipping logic
		var nX int
		var nY int
		dx := float32(p1[0] - p0[0])
		dy := float32(p1[1] - p0[1])
		if *oc&OutCodeTop != 0 {
			// Top bounds
			nY = r[1]
			nX = (*p)[0] + int(dx*(float32(r[1]-(*p)[1])/dy))
		} else if *oc&OutCodeBottom != 0 {
			// Bottom bounds
			nY = r[1] + r[3] - 1
			nX = (*p)[0] + int(dx*(float32(nY-(*p)[1])/dy))
		} else if *oc&OutCodeRight != 0 {
			// Right bounds
			nX = r[0] + r[2] - 1
			nY = (*p)[1] + int(dy*(float32(nX-(*p)[0])/dx))
		} else if *oc&OutCodeLeft != 0 {
			// Left Bounds
			nX = r[0]
			nY = (*p)[1] + int(dy*(float32(r[0]-(*p)[0])/dx))
		}
		(*p) = PointI2D{nX, nY}
		(*oc) = r.OutCode(*p)
	}
}

// UnclippedLine returns a slice of PointI2D's representing the points along a
// line from p0 to p1 in 2D space. When done with the returned slice, release it
// with ReleasePointI2DPoolSlice.
func UnclippedLine(p0, p1 PointI2D) []PointI2D {
	ret := PointI2DPool.Get()
	// Setup stepping
	dx := p1[0] - p0[0]
	dy := p1[1] - p0[1]
	stepX := 1
	if dx < 0 {
		stepX = -1
	}
	stepY := 1
	if dy < 0 {
		stepY = -1
	}
	absDx := dx
	if absDx < 0 {
		absDx = -absDx
	}
	absDy := dy
	if absDy < 0 {
		absDy = -absDy
	}
	if absDx >= absDy {
		// X-major axis
		ret = append(ret, p0)
		p := 2*absDy - absDx
		for range absDx {
			if p >= 0 {
				p0[1] += stepY
				p += 2 * (absDy - absDx)
			} else {
				p += 2 * absDy
			}
			p0[0] += stepX
			ret = append(ret, p0)
		}
	} else {
		// Y-major axis
		ret = append(ret, p0)
		p := 2*absDx - absDy
		for range absDy {
			if p >= 0 {
				p0[0] += stepX
				p += 2 * (absDx - absDy)
			} else {
				p += 2 * absDx
			}
			p0[1] += stepY
			ret = append(ret, p0)
		}
	}
	return ret
}
