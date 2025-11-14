package util

// BoundInt bounds an int to an inclusive range.
func BoundInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MaxInt returns the maximum of l and r.
func MaxInt(l, r int) int {
	if l > r {
		return l
	}
	return r
}

// MinInt returns the minimum of l and r.
func MinInt(l, r int) int {
	if l < r {
		return l
	}
	return r
}
