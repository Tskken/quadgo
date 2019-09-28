package quadgo

import (
	"fmt"
	"math"
)

// Bound is the basic rectangular bounds for nodes and entities in QuadGo.
type Bound struct {
	Min, Max, Center Point
}

// NewBound creates a new Bound struct from the given min and max points.
//
// Note: QuadGo format follows the standard 0,0 as the top left of the screen.
func NewBound(minX, minY, maxX, maxY float64) Bound {
	return Bound{
		Min: Point{X: minX, Y: minY},
		Max: Point{X: maxX, Y: maxY},
		Center: Point{
			X: maxX - ((math.Abs(maxX) - math.Abs(minX)) / 2),
			Y: maxY - ((math.Abs(maxY) - math.Abs(minY)) / 2),
		},
	}
}

// IsEqual checks if the given bound is equal to this bound.
//
// Only checks min and max points as center is based off those points
// and checking it would be redundant.
func (b Bound) IsEqual(bound Bound) bool {
	return b.Min.IsEqual(bound.Min) && b.Max.IsEqual(bound.Max)
}

// IsIntersect returns whether or not the given Bound intersects with this bound.
func (b Bound) IsIntersect(bounds Bound) bool {
	return !(bounds.Max.X < b.Min.X || bounds.Min.X > b.Max.X || bounds.Max.Y < b.Min.Y || bounds.Min.Y > b.Max.Y)
}

func (b Bound) String() string {
	return fmt.Sprintf("Min: %v, Max: %v, Center: %v\n", b.Min, b.Max, b.Center)
}
