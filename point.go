package QuadGo

import (
	"fmt"
)

// Point is the basic xy coordinates for QuadGo.
//
// QuadGo provides you with a basic implementation of Pointer through Point
type Point struct {
	x float64
	y float64
}

// NewPoint creates a new basic pointer
func NewPoint(x, y float64) Pointer {
	return &Point{x: x, y: y}
}

// X returns the x coordinate for point
func (p *Point) X() float64 {
	return p.x
}

// Y returns the y coordinate for point
func (p *Point) Y() float64 {
	return p.y
}

// XY returns the xy coordinates for point
func (p *Point) XY() (float64, float64) {
	return p.x, p.y
}

// String formats the Stringer for logging reasons
func (p *Point) String() string {
	return fmt.Sprintf("X: %f, Y: %f", p.x, p.y)
}
