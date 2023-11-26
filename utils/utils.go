package utils

func Coerce(input, min, max int) int {
	if max < min {
		newMax := min
		min = max
		max = newMax
	}
	if input < min {
		return min
	}
	if input > max {
		return max
	}
	return input
}
