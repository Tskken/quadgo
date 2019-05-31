package quadgo

import (
	"testing"
)

var (
	w              = 1024.0
	h              = 768.0
	baseBoundsList = []Bounds{
		NewBounds(0, 0, w/2, h/2),
		NewBounds(w/2, 0, w, h/2),
		NewBounds(0, h/2, w/2, h),
		NewBounds(w/2, h/2, w, h),
	}
	baseEntityList = []*Entity{
		{NewBounds(0, 0, w/2, h/2), nil},
		{NewBounds(w/2, 0, w, h/2), nil},
		{NewBounds(0, h/2, w/2, h), nil},
		{NewBounds(w/2, h/2, w, h), nil},
		{NewBounds(25, 25, 50, 50), nil},
		{NewBounds(50, 50, 100, 100), nil},
		{NewBounds(75, 75, 100, 100), nil},
		{NewBounds(10, 10, 45, 4530), nil},
		{NewBounds(15, 15, 30, 15), nil},
		{NewBounds(150, 150, 350, 350), nil},
	}
)

func TestNewQuadGo(t *testing.T) {
	t.Run("Create new QuadGo pass test", func(t *testing.T) {
		_, err := NewQuadGo(25, w, h)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Create new QuadGo maxEntities below not valid test", func(t *testing.T) {
		_, err := NewQuadGo(0, w, h)
		if err == nil {
			t.Fail()
		}
	})
}

func TestQuadGo_Insert(t *testing.T) {
	t.Run("Basic Insert Test", func(t *testing.T) {
		Quad, _ := NewQuadGo(1, w, h)

		b := NewBounds(10, 10, 60, 60)
		Quad.Insert(b)
		if Quad.entities[0].Bounds != b {
			t.Fail()
		}
	})

	t.Run("Split Test", func(t *testing.T) {
		Quad, _ := NewQuadGo(1, w, h)

		for _, b := range baseBoundsList {
			Quad.Insert(b, nil)
		}

		t.Run("Split failed tests", func(t *testing.T) {
			if len(Quad.children) == 0 {
				t.Fail()
			}
		})

		t.Run("Bottom Left Entity filled test", func(t *testing.T) {
			if len(Quad.children[bottomLeft].entities) == 0 || Quad.children[bottomLeft].entities[0].Bounds != baseBoundsList[0] {
				t.Fail()
			}
		})

		t.Run("Bottom Right Entity filled test", func(t *testing.T) {
			if len(Quad.children[bottomRight].entities) == 0 || Quad.children[bottomRight].entities[0].Bounds != baseBoundsList[1] {
				t.Fail()
			}
		})

		t.Run("Top Left Entity filled test", func(t *testing.T) {
			if len(Quad.children[topLeft].entities) == 0 || Quad.children[topLeft].entities[0].Bounds != baseBoundsList[2] {
				t.Fail()
			}
		})

		t.Run("Top Right Entity filled test", func(t *testing.T) {
			if len(Quad.children[topRight].entities) == 0 || Quad.children[topRight].entities[0].Bounds != baseBoundsList[3] {
				t.Fail()
			}
		})
	})

}

func TestQuadGo_Remove(t *testing.T) {
	Quad, _ := NewQuadGo(5, w, h)

	for _, e := range baseEntityList {
		Quad.InsertEntity(e)
	}

	t.Run("Basic Removal", func(t *testing.T) {
		Quad.Remove(baseEntityList[6])

		if Quad.IsEntity(baseEntityList[6]) {
			t.Fail()
		}
	})

	t.Run("collapse test", func(t *testing.T) {
		for _, e := range baseEntityList[4:] {
			Quad.Remove(e)
		}

		for _, c := range Quad.children {
			if len(c.children) != 0 {
				t.Fail()
			}
		}
	})

}

func TestQuadGo_Retrieve(t *testing.T) {
	Quad, _ := NewQuadGo(4, w, h)

	for _, b := range baseBoundsList {
		Quad.Insert(b)
	}

	t.Run("Basic Retrieve test", func(t *testing.T) {
		e := Quad.Retrieve(NewBounds(50, 50, 100, 100))
		if len(e) == 0 {
			t.Fail()
		}
	})

	Quad.Insert(NewBounds(400, 80, 550, 280), nil)
	Quad.Insert(NewBounds(50, 800, 70, 850), nil)

	t.Run("Retrieve only Bottom left test", func(t *testing.T) {
		e := Quad.Retrieve(NewBounds(0, 0, 50, 50))
		if len(e) > 1 || len(e) == 0 {
			t.Fail()
		}
	})
}

func TestQuadGo_IsEntity(t *testing.T) {
	Quad, _ := NewQuadGo(4, w, h)

	for _, e := range baseEntityList {
		Quad.InsertEntity(e)
	}

	t.Run("basic isEntity test", func(t *testing.T) {
		if !Quad.IsEntity(baseEntityList[4]) {
			t.Fail()
		}
	})
}

func TestQuadGo_IsIntersect(t *testing.T) {
	Quad, _ := NewQuadGo(1, w, h)

	for _, b := range baseBoundsList {
		Quad.Insert(b)
	}

	t.Run("Basic Intersect test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(20, 20, 40, 40)) {
			t.Fail()
		}
	})

	t.Run("Not Intersected test", func(t *testing.T) {
		if Quad.IsIntersect(NewBounds(-50, -50, -30, -30)) {
			t.Fail()
		}
	})

	t.Run("Intersect topLeft test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(0, Quad.bounds.width/2+20, 20, Quad.bounds.height-20)) {
			t.Fail()
		}
	})

	t.Run("Intersect TopRight test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(Quad.bounds.width/2+20, Quad.bounds.height/2+20, Quad.bounds.width/2+70, Quad.bounds.height/2+70)) {
			t.Fail()
		}
	})

	t.Run("Intersect BottomRight test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(Quad.bounds.width/2+20, 0, Quad.bounds.width/2+70, 50)) {
			t.Fail()
		}
	})
}

func TestQuadGo_Intersects(t *testing.T) {
	Quad, _ := NewQuadGo(1, w, h)

	for _, b := range baseBoundsList {
		Quad.Insert(b)
	}

	t.Run("Basic Intersects test", func(t *testing.T) {
		if len(Quad.Intersects(NewBounds(0, 0, 50, 50))) == 0 {
			t.Fail()
		}
	})
}

func BenchmarkNewQuadGo(b *testing.B) {
	b.Run("NewQuadGo pass bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := NewQuadGo(25, w, h)
			if err != nil {
				b.Fail()
			}
		}
	})

	b.Run("NewQuadGo fail bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := NewQuadGo(0, w, h)
			if err == nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkQuadGo_Insert(b *testing.B) {
	Quad, _ := NewQuadGo(1, w, h)

	for n := 0; n < b.N; n++ {
		for _, b := range baseBoundsList {
			Quad.Insert(b)
		}
	}
}

func BenchmarkQuadGo_IsIntersect(b *testing.B) {
	Quad, _ := NewQuadGo(25, w, h)

	for _, b := range baseBoundsList {
		Quad.Insert(b)
	}

	boundsf1 := NewBounds(35, 70, 85, 150)

	boundsf2 := NewBounds(-20, -50, -10, -20)

	b.Run("IsIntersect true bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if !Quad.IsIntersect(boundsf1) {
				b.Fail()
			}
		}
	})

	b.Run("IsIntersect false bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if Quad.IsIntersect(boundsf2) {
				b.Fail()
			}
		}
	})
}

func BenchmarkQuadGo_Intersects(b *testing.B) {
	Quad, _ := NewQuadGo(25, w, h)

	for _, b := range baseBoundsList {
		Quad.Insert(b)
	}

	bounds := NewBounds(25, 50, 50, 75)

	b.Run("Basic intersets bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if len(Quad.Intersects(bounds)) == 0 {
				b.Fail()
			}
		}
	})
}

func BenchmarkQuadGo_IsEntity(b *testing.B) {
	Quad, _ := NewQuadGo(4, w, h)

	for _, e := range baseEntityList {
		Quad.InsertEntity(e)
	}

	b.Run("IsEntity Pass bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if !Quad.IsEntity(baseEntityList[4]) {
				b.Fail()
			}
		}
	})

	b.Run("IsEntity fail bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if Quad.IsEntity(&Entity{}) {
				b.Fail()
			}
		}
	})
}

func BenchmarkQuadGo_Remove(b *testing.B) {
	Quad, _ := NewQuadGo(5, w, h)

	for _, e := range baseEntityList {
		Quad.InsertEntity(e)
	}

	b.Run("Remove pass bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Quad.Remove(baseEntityList[6])
		}
	})

	b.Run("Collapse pass bench", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			for _, e := range baseEntityList[4:] {
				Quad.Remove(e)
			}
		}
	})
}
