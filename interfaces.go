package QuadGo

type Bounder interface {
	Center() Pointer
	Bounds() (min, max Pointer)
	Min() Pointer
	Max() Pointer
	W() float64
	H() float64
}

type Pointer interface {
	X() float64
	Y() float64
	XY() (float64, float64)
}

type Entity interface {
	IsIntersect(Entity) bool
	Bounder
}