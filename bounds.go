package QuadGo

import (
	"math"
)

// Bounds is the basic AABB or rect bounds structure for QuadGo.
type Bounds struct {
	min, max      point
	width, height float64
}

// NewBounds creates a new Bounds struct from the given min and max x y coordinates.
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

// Center returns the center x y coordinates of the bounds as a point.
func (b *Bounds) center() point {
	return point{x: b.max.x - (b.width / 2), y: b.max.y - (b.height / 2)}
}

// isIntersect returns whether or not the given bounds intersects with bounds.
func (b *Bounds) isIntersect(bounds Bounds) bool {
	// check if given Entity does not fit with in node bounds
	if bounds.max.x < b.min.x || bounds.min.x > b.max.x || bounds.max.y < b.min.y || bounds.min.y > b.max.y {
		return false
	}

	return true
}
