package core

import "math"

// NegMod : Extends modulus operator for
// negative numbers such that the results
// makes some damn sense
func NegMod(m, n int32) int32 {

	return (m%n + n) % n
}

// MinInt32 : Minimum of two signed integers
func MinInt32(x, y int32) int32 {

	if x < y {

		return x
	}
	return y
}

// MinUInt32 : Minimum of two unsigned integers
func MinUInt32(x, y uint32) uint32 {

	if x < y {

		return x
	}
	return y
}

// MaxInt32 : Maximum of two signed integers
func MaxInt32(x, y int32) int32 {

	if x > y {

		return x
	}
	return y
}

// MaxInt32InSlice : Return the maximal element in
// a slice
func MaxInt32InSlice(arr []int32) int32 {

	max := arr[0]
	for _, m := range arr {

		if m > max {

			max = m
		}
	}
	return max
}

// MaxUInt32 : Maximum of two unsigned integers
func MaxUInt32(x, y uint32) uint32 {

	if x > y {

		return x
	}
	return y
}

// RoundFloat32 : Round a 32-bit floating point
// number
func RoundFloat32(x float32) int32 {

	return RoundFloat64(float64(x))
}

// RoundFloat64 : Round a 64-bit floating point
// number
func RoundFloat64(x float64) int32 {

	return int32(math.Round(x))
}

// ClampFloat32 : "Clamps" the given number to the interval
// [min, max]
func ClampFloat32(x float32, min float32, max float32) float32 {

	return float32(math.Min(float64(max),
		math.Max(float64(x), float64(min))))
}

// ClampInt32 : "Clamps" the given integer to the interval
// [min, max]
func ClampInt32(x int32, min int32, max int32) int32 {

	return int32(MinInt32(max, MaxInt32(x, min)))
}

// HypotInt32 : Like math.Hypot, but works with integers
func HypotInt32(x, y int32) int32 {

	return RoundFloat64(math.Hypot(float64(x), float64(y)))
}
