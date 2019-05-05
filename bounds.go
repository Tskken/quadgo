package QuadGo

import (
	"fmt"
	"math"
)

// Bounds is the basic entity and bounds structure for QuadGo.
//
// QuadGo provides a very basic bounds and entity structure for basic collision detection.
// You can creat your own entity and bounds by implementing the Bounder, Pointer, and Entity interfaces.
type Bounds struct {
	// Contains unexported fields
	min    Pointer
	max    Pointer
	width  float64
	height float64
}

// NewBounds creates a new Bounds struct from the given min x y coordinate and the width and height.
//
// Note: QuadGo minX and minY points are the bottom left point of the bounding rectangle.
func NewBounds(minX, minY, width, height float64) Bounder {
	return &Bounds{
		min:    &Point{x: minX, y: minY},
		max:    &Point{x: minX + width, y: minY + height},
		width:  width,
		height: height,
	}
}

// ToBounds creates a new Bounds struct from the given min and max x y coordinates.
//
// Note: QuadGo format has minX and minY as the bottom left and maxX and maxY as top right.
func ToBounds(minX, minY, maxX, maxY float64) Bounder {
	return &Bounds{
		min:    &Point{x: minX, y: minY},
		max:    &Point{x: maxX, y: maxY},
		width:  math.Abs(maxX) - math.Abs(minX),
		height: math.Abs(maxY) - math.Abs(minY),
	}
}

// Center returns the center Point of the bounds.
func (b *Bounds) Center() Pointer {
	return &Point{
		x: b.min.X() + (b.width / 2),
		y: b.min.Y() + (b.height / 2),
	}
}

// IsIntersect returns weather or not the given entity intersects with the bounds.
func (b *Bounds) IsIntersect(entity Entity) bool {
	min, max := entity.Bounds()
	// Left of entity
	if max.X() < b.min.X() {
		return false
	}
	// Right of entity
	if min.X() > b.max.X() {
		return false
	}
	// Above entity
	if max.Y() < b.min.Y() {
		return false
	}
	// Below entity
	if min.Y() > b.max.Y() {
		return false
	}

	return true
}

// Bounds returns the min and max xy coordinates of Bounds.
func (b *Bounds) Bounds() (Pointer, Pointer) {
	return b.min, b.max
}

// Min returns the min xy coordinates of bounds.
func (b *Bounds) Min() Pointer {
	return b.min
}

// Max returns the max xy coordinates of bounds.
func (b *Bounds) Max() Pointer {
	return b.max
}

// W returns the width of bounds.
func (b *Bounds) W() float64 {
	return b.width
}

// H returns the height of bounds.
func (b *Bounds) H() float64 {
	return b.height
}

// String formats the Stringer for logging reasons.
func (b *Bounds) String() string {
	return fmt.Sprintf("Min: %v, Max: %v, Width: %f, Height: %f", b.min, b.max, b.width, b.height)
}
