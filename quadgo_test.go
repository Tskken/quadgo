package QuadGo

import (
	"testing"
)

var (
	RootBounds = NewBounds(0, 0, 1024, 768)
	Quad       *QuadGo
)

func TestQuadGo_Insert(t *testing.T) {
	t.Run("Basic Insert Test", func(t *testing.T) {
		Quad = NewQuadGo(1, RootBounds)

		entity := NewBounds(10, 10, 50, 50)
		Quad.Insert(entity)
		if Quad.root.entities[0] != entity {
			t.Fail()
		}
	})

	t.Run("Split Test", func(t *testing.T) {
		Quad = NewQuadGo(1, RootBounds)

		entities := []Entity{
			NewBounds(0, 0, 50, 50),
			NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
		}

		for _, e := range entities {
			Quad.Insert(e)
		}

		t.Run("Split failed tests", func(t *testing.T) {
			if len(Quad.root.children) == 0 {
				t.Fail()
			}
		})

		t.Run("Bottom Left entity filled test", func(t *testing.T) {
			if len(Quad.root.children[bottomLeft].entities) == 0 || Quad.root.children[bottomLeft].entities[0] != entities[0] {
				t.Fail()
			}
		})

		t.Run("Bottom Right entity filled test", func(t *testing.T) {
			if len(Quad.root.children[bottomRight].entities) == 0 || Quad.root.children[bottomRight].entities[0] != entities[2] {
				t.Fail()
			}
		})

		t.Run("Top Left entity filled test", func(t *testing.T) {
			if len(Quad.root.children[topLeft].entities) == 0 || Quad.root.children[topLeft].entities[0] != entities[1] {
				t.Fail()
			}
		})

		t.Run("Top Right entity filled test", func(t *testing.T) {
			if len(Quad.root.children[topRight].entities) == 0 || Quad.root.children[topRight].entities[0] != entities[3] {
				t.Fail()
			}
		})
	})

}

func TestQuadGo_Remove(t *testing.T) {
	Quad = NewQuadGo(5, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(25, 25, 25, 25),
		NewBounds(50, 50, 50, 50),
		NewBounds(75, 75, 25, 25),
		NewBounds(10, 10, 35, 35),
		NewBounds(15, 15, 15, 15),
		NewBounds(150, 150, 200, 200),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	t.Run("Basic Removal", func(t *testing.T) {
		Quad.Remove(entities[6])

		if Quad.IsEntity(entities[5]) {
			t.Fail()
		}
	})

}

func TestQuadGo_Retrieve(t *testing.T) {
	Quad = NewQuadGo(4, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	t.Run("Basic Retrieve test", func(t *testing.T) {
		e := Quad.Retrieve(NewBounds(50, 50, 50, 50))
		if len(e) == 0 {
			t.Fail()
		}
	})

	Quad.Insert(NewBounds(400, 80, 150, 200))
	Quad.Insert(NewBounds(50, 800, 20, 50))

	t.Run("Retrieve only Bottom left test", func(t *testing.T) {
		e := Quad.Retrieve(NewBounds(0, 0, 50, 50))
		if len(e) > 1 || len(e) == 0 {
			t.Fail()
		}
	})
}

func TestQuadGo_Find(t *testing.T) {
	Quad = NewQuadGo(4, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	t.Run("basic isEntity test", func(t *testing.T) {
		if !Quad.IsEntity(entities[0]) {
			t.Fail()
		}
	})
}

func TestQuadGo_IsIntersect(t *testing.T) {
	Quad = NewQuadGo(1, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	t.Run("Basic Intersect test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(20, 20, 20, 20)) {
			t.Fail()
		}
	})

	t.Run("Not Intersected test", func(t *testing.T) {
		if Quad.IsIntersect(NewBounds(-50, -50, 20, 20)) {
			t.Fail()
		}
	})

	t.Run("Intersect topLeft test", func(t *testing.T) {
		if !Quad.IsIntersect(ToBounds(0, Quad.root.bounds.H()/2+20, 20, Quad.root.bounds.H()-20)) {
			t.Fail()
		}
	})

	t.Run("Intersect TopRight test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(Quad.root.bounds.W()/2+20, Quad.root.bounds.H()/2+20, 50, 50)) {
			t.Fail()
		}
	})

	t.Run("Intersect BottomRight test", func(t *testing.T) {
		if !Quad.IsIntersect(NewBounds(Quad.root.bounds.W()/2+20, 0, 50, 50)) {
			t.Fail()
		}
	})
}

func TestQuadGo_Intersects(t *testing.T) {
	Quad = NewQuadGo(1, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	t.Run("Basic Intersects test", func(t *testing.T) {
		if len(Quad.Intersects(NewBounds(0, 0, 50, 50))) == 0 {
			t.Fail()
		}
	})
}

func BenchmarkNewQuadGo(b *testing.B) {
	Quad = NewQuadGo(25, RootBounds)
}

func BenchmarkQuadGo_Insert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Quad = NewQuadGo(1, RootBounds)

		entities := []Entity{
			NewBounds(0, 0, 50, 50),
			NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
		}

		for _, e := range entities {
			Quad.Insert(e)
		}
	}
}

func BenchmarkQuadGo_IsIntersect(b *testing.B) {
	Quad = NewQuadGo(25, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(150, 40, 300, 40),
		NewBounds(34, 65, 234, 680),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	entity := NewBounds(35, 70, 50, 80)

	for n := 0; n < b.N; n++ {
		if !Quad.IsIntersect(entity) {
			b.Fail()
		}
	}
}

func BenchmarkQuadGo_Find(b *testing.B) {
	Quad = NewQuadGo(4, RootBounds)

	entities := []Entity{
		NewBounds(0, 0, 50, 50),
		NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
		NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
	}

	for _, e := range entities {
		Quad.Insert(e)
	}

	for n := 0; n < b.N; n++ {
		if !Quad.IsEntity(entities[0]) {
			b.Fail()
		}
	}
}

func BenchmarkQuadGo_Remove(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Quad = NewQuadGo(5, RootBounds)

		entities := []Entity{
			NewBounds(0, 0, 50, 50),
			NewBounds(0, Quad.root.bounds.H()/2+50, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, 0, 50, 50),
			NewBounds(Quad.root.bounds.W()/2+50, Quad.root.bounds.H()/2+50, 50, 50),
			NewBounds(25, 25, 25, 25),
			NewBounds(50, 50, 50, 50),
			NewBounds(75, 75, 25, 25),
			NewBounds(10, 10, 35, 35),
			NewBounds(15, 15, 15, 15),
			NewBounds(150, 150, 200, 200),
		}

		for _, e := range entities {
			Quad.Insert(e)
		}

		Quad.Remove(entities[6])
	}
}
