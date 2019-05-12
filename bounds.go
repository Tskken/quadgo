package QuadGo

import (
	"math"
)

// Bounds is the basic Entity and bounds structure for QuadGo.
//
// QuadGo provides a very basic bounds and Entity structure for basic collision detection.
// You can creat your own Entity and bounds by implementing the Bounder and Entity interfaces.
type Bounds struct {
	// Contains unexported fields
	min, max      point
	width, height float64
}

// ToBounds creates a new Bounds struct from the given min and max x y coordinates.
//
// Note: QuadGo format has minX and minY as the bottom left and maxX and maxY as top right.
func NewBounds(minX, minY, maxX, maxY float64) Bounds {
	return Bounds{
		min:    point{x: minX, y: minY},
		max:    point{x: maxX, y: maxY},
		width:  math.Abs(maxX) - math.Abs(minX),
		height: math.Abs(maxY) - math.Abs(minY),
	}
}

// Center returns the center x y coordinates of the bounds.
func (b *Bounds) center() point {
	return point{x: b.max.x - (b.width / 2), y: b.max.y - (b.height / 2)}
}

// isIntersect returns whether or not the given Entity intersects with the bounds.
func (b *Bounds) isIntersect(bounds Bounds) bool {
	// check if given Entity does not fit with in node bounds
	if bounds.max.x < b.min.x || bounds.min.x > b.max.x || bounds.max.y < b.min.y || bounds.min.y > b.max.y {
		return false
	}

	return true
}

//func (b *Bounds) validate() bool {
//	return b.min.x >= b.max.x || b.min.y >= b.max.y
//}

//// String formats the Stringer for logging reasons.
//func (b *Bounds) String() string {
//	return fmt.Sprintf("Min: %f, %f, Max: %f, %f, Width: %f, Height: %f", b.minX, b.minY, b.maxX, b.maxY, b.width, b.height)
//}
