package lss

import "math"

func IAbs(x int) int {
	if x >= 0 {
		return x
	}
	return x * -1
}

func NumDigits(x int) int {
	return int(math.Floor(math.Log10(float64(IAbs(x))))) + 1
}
