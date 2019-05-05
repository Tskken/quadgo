package QuadGo

import "math"

type Bounds struct {
	min    Pointer
	max    Pointer
	width  float64
	height float64
}

func NewBounds(minX, minY, width, height float64) Bounder {
	return &Bounds{
		min:    &point{x: minX, y: minY},
		max:    &point{x: minX + width, y: minY + height},
		width:  width,
		height: height,
	}
}

func ToBounds(minX, minY, maxX, maxY float64) Bounder {
	return &Bounds{
		min:    &point{x: minX, y: minY},
		max:    &point{x: maxX, y: maxY},
		width:  math.Abs(maxX) - math.Abs(minX),
		height: math.Abs(maxY) - math.Abs(minY),
	}
}

func (b *Bounds) Center() Pointer {
	return &point{
		x: b.min.X() + (b.width / 2),
		y: b.min.Y() + (b.height / 2),
	}
}

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

func (b *Bounds) Bounds() (Pointer, Pointer) {
	return b.min, b.max
}

func (b *Bounds) Min() Pointer {
	return b.min
}

func (b *Bounds) Max() Pointer {
	return b.max
}

func (b *Bounds) W() float64 {
	return b.width
}

func (b *Bounds) H() float64 {
	return b.height
}
