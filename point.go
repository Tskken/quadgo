package quadgo

import "fmt"

// Point is the basic 2D coordinate structure for QuadGo
type Point struct {
	X, Y float64
}

// NewPoint creates a new point for the given x and y positions.
func NewPoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// IsEqual checks if the given point is equal to this point.
func (p Point) IsEqual(point Point) bool {
	return p.X == point.X && p.Y == point.Y
}

func (p Point) String() string {
	return fmt.Sprintf("X: %v, Y: %v", p.X, p.Y)
}
