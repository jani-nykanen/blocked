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

	return int32(math.Round(float64(x)))
}
