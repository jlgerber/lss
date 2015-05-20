package lss

import "math"

// IAbs - Integer Abs call. Like math.Abs, but for integers
func IAbs(x int) int {
	if x >= 0 {
		return x
	}
	return x * -1
}

// NumDigits - Given an integer, return the number of digits of which it is comprised. For example,
// NumDigits(11) => 2 NumDigits(3) => 1
func NumDigits(x int) int {
	return int(math.Floor(math.Log10(float64(IAbs(x))))) + 1
}
