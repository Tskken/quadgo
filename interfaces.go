package QuadGo

// Bounder interface for implementing all bounds for the quadtree.
// Any entity that you want to add to the quadtree must implement this interface.
type Bounder interface {
	IsIntersect(Entity) bool
	Center() (x, y float64)
	Bounds() (minX, minY, maxX, maxY float64)
	Min() (x, y float64)
	Max() (x, y float64)
	W() float64
	H() float64
}

// Entity extends the implementation of the Bounds interface.
// Any entity you want to add to the quadtree must implement this interface.
//
// Note: This currently is exactly the same as Bounder but may change in the future if
// I find any need for Entity specific functions. For know it is just hear to make reading
// the code easier as having Bounder in areas were you were dealing with entities seemed odd.
type Entity interface {
	Bounder
}
