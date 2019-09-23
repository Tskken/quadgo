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
// Note: QuadGo format fallows the default 0,0 as the top left of the screen.
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
func (b Bound) IsEqual(bound Bound) bool {
	return b.Min.IsEqual(bound.Min) && b.Max.IsEqual(bound.Max)
}

// IsIntersectBound returns whether or not the given Bound intersects with this bound.
func (b Bound) IsIntersectBound(bounds Bound) bool {
	return !(bounds.Max.X < b.Min.X || bounds.Min.X > b.Max.X || bounds.Max.Y < b.Min.Y || bounds.Min.Y > b.Max.Y)
}

// IsIntersectPoint returns whether or not the given point intersects with this bound.
func (b Bound) IsIntersectPoint(point Point) bool {
	return !(point.X < b.Min.X || point.X > b.Max.X || point.Y < b.Min.Y || point.Y > b.Max.Y)
}

func (b Bound) String() string {
	return fmt.Sprintf("Min: %v, Max: %v, Center: %v\n", b.Min, b.Max, b.Center)
}
