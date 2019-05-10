package QuadGo

import (
	"fmt"
	"math"
)

// Bounds is the basic entity and bounds structure for QuadGo.
//
// QuadGo provides a very basic bounds and entity structure for basic collision detection.
// You can creat your own entity and bounds by implementing the Bounder and Entity interfaces.
type Bounds struct {
	// Contains unexported fields
	minX, minY float64
	maxX, maxY float64
	width, height float64
}

// NewBounds creates a new Bounds struct from the given min x y coordinate and the width and height.
//
// Note: QuadGo minX and minY points are the bottom left point of the bounding rectangle.
func NewBounds(minX, minY, width, height float64) *Bounds {
	return &Bounds{
		minX: minX, minY: minY,
		maxX: minX + width, maxY: minY + height,
		width:  width, height: height,
	}
}

// ToBounds creates a new Bounds struct from the given min and max x y coordinates.
//
// Note: QuadGo format has minX and minY as the bottom left and maxX and maxY as top right.
func ToBounds(minX, minY, maxX, maxY float64) *Bounds {
	return &Bounds{
		minX: minX, minY: minY,
		maxX: maxX, maxY: maxY,
		width:  math.Abs(maxX) - math.Abs(minX),
		height: math.Abs(maxY) - math.Abs(minY),
	}
}

// Center returns the center x y coordinates of the bounds.
func (b *Bounds) Center() (x, y float64) {
	return b.minX + (b.width / 2), b.minY + (b.height / 2)
}

// IsIntersect returns whether or not the given entity intersects with the bounds.
func (b *Bounds) IsIntersect(entity Entity) bool {
	minX, minY, maxX, maxY := entity.Bounds()

	// Left of entity
	if maxX < b.minX || minX > b.maxX || maxY < b.minY || minY > b.maxY {
		return false
	}

	return true
}

// Bounds returns the min and max xy coordinates of Bounds.
func (b *Bounds) Bounds() (minX, minY float64, maxX, maxY float64) {
	return b.minX, b.minY, b.maxX, b.maxY
}

// Min returns the min xy coordinates of bounds.
func (b *Bounds) Min() (x, y float64) {
	return b.minX, b.minY
}

// Max returns the max xy coordinates of bounds.
func (b *Bounds) Max() (x, y float64) {
	return b.maxX, b.maxY
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
	return fmt.Sprintf("Min: %f, %f, Max: %f, %f, Width: %f, Height: %f", b.minX, b.minY, b.maxX, b.maxY, b.width, b.height)
}
