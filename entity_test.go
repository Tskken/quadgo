package quadgo

import "testing"

func TestNewEntity(t *testing.T) {
	t.Run("basic create new entity test", func(t *testing.T) {
		e := NewEntity(0, 0, 50, 50)
		b := NewBound(0, 0, 50, 50)

		if e.Bound != b {
			t.Errorf("entity min and max bounds not set correctly. Bounds: %v", e.Bound)
		}

		if e.Objects != nil {
			t.Errorf("objects were not defaulted to nil")
		}
	})
}

func TestNewBound(t *testing.T) {
	t.Run("Basic create new bound test", func(t *testing.T) {
		b := NewBound(0, 0, 50, 50)

		if b.Width != 50 && b.Height != 50 {
			t.Errorf("width and height were not set correctly. Width and height are %v, %v", b.Width, b.Height)
		}

		if b.Min.X != 0 && b.Min.Y != 0 && b.Max.X != 50 && b.Max.Y != 50 {
			t.Errorf("min and max points were not set correctly. Min and max points are %v, %v", b.Min, b.Max)
		}

		if b.Center.X != 25 && b.Center.Y != 25 {
			t.Errorf("center was not set correctly. Center is %v", b.Center)
		}
	})
}

func TestBound_IsIntersectPoint(t *testing.T) {
	t.Run("point intersect test", func(t *testing.T) {
		b := NewBound(0, 0, 50, 50)

		if !b.IsIntersectPoint(Point{25, 25}) {
			t.Errorf("IsIntersectPoint did not correctly calculate intersect")
		}
	})
}

func TestBound_IsIntersectBound(t *testing.T) {
	t.Run("bound intersect test", func(t *testing.T) {
		b := NewBound(0, 0, 50, 50)

		if !b.IsIntersectBound(NewBound(5, 5, 20, 20)) {
			t.Errorf("IsInetersectBound did not correctly calcuate intersect")
		}
	})
}
