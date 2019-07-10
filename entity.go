package quadgo

import (
	"fmt"
	"math"
)

// Entities is a list of entities.
type Entities []*Entity

// Entity is the basic Entity stricture type for QuadGo.
//
// Entity holds the Bound information for an entity in the tree and also a list of interface{} which can hold
// any data that you would want to store in the entity.
type Entity struct {
	Bound

	Objects []interface{}
}

// NewEntity creates a new entity from the given min and max points and any given objects.
//
// The given objects can be any data that you want to hold with in the entity for the given bounds.
func NewEntity(minX, minY, maxX, maxY float64, objs ...interface{}) *Entity {
	return &Entity{
		Bound:   NewBound(minX, minY, maxX, maxY),
		Objects: objs,
	}
}

func (e *Entity) String() string {
	return fmt.Sprintf("Bounds: %v\n Objects: %v\n", e.Bound, e.Objects)
}

// Bound is the basic rectangular bounds for nodes and entities in QuadGo.
type Bound struct {
	Min, Max, Center Point
	Width, Height    float64
}

var ZB = Bound{}

// NewBound creates a new Bound struct from the given min and max points.
//
// Note: QuadGo format has min as the bottom left and max as top right.
func NewBound(minX, minY, maxX, maxY float64) Bound {
	w := math.Abs(maxX) - math.Abs(minX)
	h := math.Abs(maxY) - math.Abs(minY)
	return Bound{
		Min:    Point{X: minX, Y: minY},
		Max:    Point{X: maxX, Y: maxY},
		Center: Point{X: maxX - (w / 2), Y: maxY - (h / 2)},
		Width:  w,
		Height: h,
	}
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
	return fmt.Sprintf("Min: %v, Max: %v, Center: %v\n Width: %v, Height %v\n", b.Min, b.Max, b.Center, b.Width, b.Height)
}

// Point is the basic X Y coordinate structure for QuadGo
type Point struct {
	X, Y float64
}

var ZP = Point{}

func (p Point) String() string {
	return fmt.Sprintf("X: %v, Y: %v", p.X, p.Y)
}
