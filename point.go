package QuadGo

type point struct {
	x float64
	y float64
}

func (p *point) X() float64 {
	return p.x
}

func (p *point) Y() float64 {
	return p.y
}

func (p *point) XY() (float64, float64) {
	return p.x, p.y
}
