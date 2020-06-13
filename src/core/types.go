package core

import "math"

// Vector2 : a two component vector with float point
// components
type Vector2 struct {
	X, Y float32
}

// NewVector2 : Constructs a new vector2
func NewVector2(x, y float32) Vector2 {

	return Vector2{X: x, Y: y}
}

// Normalize : Normalizes the vector, returns a copy
func (v *Vector2) Normalize() Vector2 {

	const EPSILON = 0.0001

	len := v.Length()
	if len < EPSILON {

		v.X = 0
		v.Y = 0
	} else {

		v.X /= len
		v.Y /= len
	}

	return NewVector2(v.X, v.Y)
}

// Length : Returns length of the vector
func (v *Vector2) Length() float32 {

	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Point : A 2-components vector, integer components
type Point struct {
	X, Y int32
}

// NewPoint : Constructor for point
func NewPoint(x, y int32) Point {

	return Point{X: x, Y: y}
}

// Rectangle : A 4-component vector, integer components
type Rectangle struct {
	X, Y, W, H int32
}

// NewRect : Constructor for rectangle
func NewRect(x, y, w, h int32) Rectangle {

	return Rectangle{X: x, Y: y, W: w, H: h}
}

// Color : RGBA color
type Color struct {
	R, G, B, A uint8
}

// NewRGBA : Constructor for color, all components
func NewRGBA(r, g, b, a uint8) Color {

	return Color{R: r, G: g, B: b, A: a}
}

// NewRGB : Constructor for color, no alpha component
func NewRGB(r, g, b uint8) Color {

	return Color{R: r, G: g, B: b, A: 255}
}
