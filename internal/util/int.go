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
