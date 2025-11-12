package util

// MaxFloat32 returns the maximum value passed to the function.
func MaxFloat32(args ...float32) float32 {
	if len(args) == 0 {
		return 0
	}
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

// MinFloat32 returns the minimum value passed to the function.
func MinFloat32(args ...float32) float32 {
	if len(args) == 0 {
		return 0
	}
	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}
	return min
}
