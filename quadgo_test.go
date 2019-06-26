package quadgo

import (
	"testing"
)

/*
	TODO: Redo testing to be better and require less maintenance when changes to the library are made.
*/

func TestNew(t *testing.T) {
	t.Run("basic new test, no options given", func(t *testing.T) {
		qTree := New()
		if qTree.bounds.Width != defaultOption.Width {
			t.Errorf("defualt tree width was not set to %v, insted it was set to %v", defaultOption.Width, qTree.bounds.Width)
		}

		if qTree.bounds.Height != 768 {
			t.Errorf("default tree height was not set to %v it was set insted to %v", defaultOption.Height, qTree.bounds.Height)
		}

		if qTree.entities == nil {
			t.Error("entities were not initialized")
		}

		if qTree.children == nil {
			t.Error("children were not initialized")
		}

		if qTree.maxDepth != defaultOption.MaxDepth {
			t.Errorf("max depth default was not set to %v, instead it was set to %v", defaultOption.MaxDepth, qTree.maxDepth)
		}

		if cap(qTree.entities) != defaultOption.MaxEntities {
			t.Errorf("max entities default was not set to %v, insted it was set to %v", defaultOption.MaxEntities, cap(qTree.entities))
		}
	})

	t.Run("new test with SetBound option", func(t *testing.T) {
		qTree := New(SetBounds(1920, 1080))

		if qTree.bounds.Width != 1920 {
			t.Errorf("SetBounds did not set new bounds for width. The width for the tree still are %v", qTree.bounds.Width)
		}

		if qTree.bounds.Height != 1080 {
			t.Errorf("SetBounds did not set new bounds for height, The height for the tree still are %v", qTree.bounds.Height)
		}
	})

	t.Run("new test for SetMaxDepth option", func(t *testing.T) {
		qTree := New(SetMaxDepth(4))

		if qTree.maxDepth != 4 {
			t.Errorf("SetMaxDepth did not set the new max depth for the tree. The max depth is still %v", qTree.maxDepth)
		}
	})

	t.Run("new test for SetMaxEntities option", func(t *testing.T) {
		qTree := New(SetMaxEntities(5))

		if cap(qTree.entities) != 5 {
			t.Errorf("SetMaxEntities did not set the trees new max entitie value. Its max entities is still %v", cap(qTree.entities))
		}
	})
}

func TestQuadGo_Insert(t *testing.T) {
	t.Run("basic insert test", func(t *testing.T) {
		qTree := New()

		qTree.Insert(0, 0, 50, 50)

		if len(qTree.entities) != 1 {
			t.Errorf("no entity was inserted in to tree. Entities count is %v", len(qTree.entities))
		}
	})

	t.Run("split on insert test", func(t *testing.T) {
		qTree := New(SetMaxEntities(2))

		blCenter := NewBound(0, 0, qTree.bounds.Center.X, qTree.bounds.Center.Y).Center
		brCenter := NewBound(qTree.bounds.Center.X, 0, qTree.bounds.Width, qTree.bounds.Center.Y).Center
		tlCenter := NewBound(0, qTree.bounds.Center.Y, qTree.bounds.Center.X, qTree.bounds.Height).Center
		trCenter := NewBound(qTree.bounds.Center.X, qTree.bounds.Center.Y, qTree.bounds.Width, qTree.bounds.Height).Center

		qTree.Insert(blCenter.X-25, blCenter.Y-25, blCenter.X+25, blCenter.Y+25)
		qTree.Insert(brCenter.X-25, brCenter.Y-25, brCenter.X+25, brCenter.Y+25)
		qTree.Insert(tlCenter.X-25, tlCenter.Y-25, tlCenter.X+25, tlCenter.Y+25)
		qTree.Insert(trCenter.X-25, trCenter.Y-25, trCenter.X+25, trCenter.Y+25)

		if len(qTree.children) != 4 {
			t.Errorf("node did not split on insert")
		}

		if len(qTree.entities) != 0 {
			t.Errorf("entities not removed from parent node")
		}

		eCount := 0

		for _, c := range qTree.children {
			eCount += len(c.entities)
		}

		if eCount != 4 {
			t.Errorf("entities were not all moved to there children nodes. Entity count is %v", eCount)
		}
	})

	t.Run("split to correct quadrants test", func(t *testing.T) {
		qTree := New(SetMaxEntities(1))

		blCenter := NewBound(0, 0, qTree.bounds.Center.X, qTree.bounds.Center.Y).Center
		brCenter := NewBound(qTree.bounds.Center.X, 0, qTree.bounds.Width, qTree.bounds.Center.Y).Center
		tlCenter := NewBound(0, qTree.bounds.Center.Y, qTree.bounds.Center.X, qTree.bounds.Height).Center
		trCenter := NewBound(qTree.bounds.Center.X, qTree.bounds.Center.Y, qTree.bounds.Width, qTree.bounds.Height).Center

		qTree.Insert(blCenter.X-25, blCenter.Y-25, blCenter.X+25, blCenter.Y+25)
		qTree.Insert(brCenter.X-25, brCenter.Y-25, brCenter.X+25, brCenter.Y+25)
		qTree.Insert(tlCenter.X-25, tlCenter.Y-25, tlCenter.X+25, tlCenter.Y+25)
		qTree.Insert(trCenter.X-25, trCenter.Y-25, trCenter.X+25, trCenter.Y+25)

		if len(qTree.children[bottomLeft].entities) != 1 &&
			len(qTree.children[bottomRight].entities) != 1 &&
			len(qTree.children[topLeft].entities) != 1 &&
			len(qTree.children[topRight].entities) != 1 {
			t.Errorf("entities were not split correctly")
		}
	})

	t.Run("hit max depth test", func(t *testing.T) {
		qTree := New(
			SetMaxEntities(1),
			SetMaxDepth(0),
		)

		blCenter := NewBound(0, 0, qTree.bounds.Center.X, qTree.bounds.Center.Y).Center
		brCenter := NewBound(qTree.bounds.Center.X, 0, qTree.bounds.Width, qTree.bounds.Center.Y).Center
		tlCenter := NewBound(0, qTree.bounds.Center.Y, qTree.bounds.Center.X, qTree.bounds.Height).Center
		trCenter := NewBound(qTree.bounds.Center.X, qTree.bounds.Center.Y, qTree.bounds.Width, qTree.bounds.Height).Center

		qTree.Insert(blCenter.X-25, blCenter.Y-25, blCenter.X+25, blCenter.Y+25)
		qTree.Insert(brCenter.X-25, brCenter.Y-25, brCenter.X+25, brCenter.Y+25)
		qTree.Insert(tlCenter.X-25, tlCenter.Y-25, tlCenter.X+25, tlCenter.Y+25)
		qTree.Insert(trCenter.X-25, trCenter.Y-25, trCenter.X+25, trCenter.Y+25)

		if len(qTree.children) != 0 {
			t.Errorf("split went past max depth of 1")
		}

		if len(qTree.entities) != 4 {
			t.Errorf("no entities were added. Failed insert with depth check")
		}
	})
}

func TestQuadGo_InsertEntity(t *testing.T) {
	t.Run("basic insert entity test", func(t *testing.T) {
		qTree := New()

		e := NewEntity(0, 0, 50, 50)

		qTree.InsertEntity(e)

		if len(qTree.entities) != 1 {
			t.Errorf("entity was not added to tree")
		}
	})
}

//var (
//	w              = 1024.0
//	h              = 768.0
//	baseBoundsList = []Bound{
//		NewBound(ZP, Point{w/2, h/2}),
//		NewBound(Point{w/2, 0}, Point{w, h/2}),
//		NewBound(Point{0, h/2}, Point{w/2, h}),
//		NewBound(Point{w/2, h/2}, Point{w, h}),
//	}
//	baseEntityList = Entities{
//		{NewBound(ZP, Point{w/2, h/2}), nil},
//		{NewBound(Point{w/2, 0}, Point{w, h/2}), nil},
//		{NewBound(Point{0, h/2}, Point{w/2, h}), nil},
//		{NewBound(Point{w/2, h/2}, Point{w, h}), nil},
//		{NewBound(Point{25, 25}, Point{50, 50}), nil},
//		{NewBound(Point{50, 50}, Point{100, 100}), nil},
//		{NewBound(Point{75, 75}, Point{100, 100}), nil},
//		{NewBound(Point{10, 10}, Point{45, 4530}), nil},
//		{NewBound(Point{15, 15}, Point{30, 15}), nil},
//		{NewBound(Point{150, 150}, Point{350, 350}), nil},
//	}
//)
//
//func TestNewQuadGo(t *testing.T) {
//	t.Run("Create new QuadGo pass test", func(t *testing.T) {
//		q := New(w, h, 25)
//		if q == nil {
//			t.Error("failed to create quadgo.New()")
//		}
//	})
//}
//
//func TestQuadGo_Insert(t *testing.T) {
//	t.Run("Basic Insert Test", func(t *testing.T) {
//		Quad := New(w, h, 1)
//
//		b := NewBound(Point{10, 10}, Point{60, 60})
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//		if Quad.entities[0].Bound != b {
//			t.Fail()
//		}
//	})
//
//	t.Run("Split Test", func(t *testing.T) {
//		Quad := New(w, h, 1)
//
//		for _, b := range baseBoundsList {
//			Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//		}
//
//		t.Run("Split failed tests", func(t *testing.T) {
//			if len(Quad.children) == 0 {
//				t.Fail()
//			}
//		})
//
//		t.Run("Bottom Left Entity filled test", func(t *testing.T) {
//			if len(Quad.children[bottomLeft].entities) == 0 || Quad.children[bottomLeft].entities[0].Bound != baseBoundsList[0] {
//				t.Fail()
//			}
//		})
//
//		t.Run("Bottom Right Entity filled test", func(t *testing.T) {
//			if len(Quad.children[bottomRight].entities) == 0 || Quad.children[bottomRight].entities[0].Bound != baseBoundsList[1] {
//				t.Fail()
//			}
//		})
//
//		t.Run("Top Left Entity filled test", func(t *testing.T) {
//			if len(Quad.children[topLeft].entities) == 0 || Quad.children[topLeft].entities[0].Bound != baseBoundsList[2] {
//				t.Fail()
//			}
//		})
//
//		t.Run("Top Right Entity filled test", func(t *testing.T) {
//			if len(Quad.children[topRight].entities) == 0 || Quad.children[topRight].entities[0].Bound != baseBoundsList[3] {
//				t.Fail()
//			}
//		})
//	})
//
//}
//
//func TestQuadGo_Remove(t *testing.T) {
//	Quad := New(w, h, 5)
//
//	for _, e := range baseEntityList {
//		Quad.InsertEntity(e)
//	}
//
//	t.Run("Basic Removal", func(t *testing.T) {
//		Quad.Remove(baseEntityList[6])
//
//		if Quad.IsEntity(baseEntityList[6]) {
//			t.Fail()
//		}
//	})
//
//	t.Run("collapse test", func(t *testing.T) {
//		for _, e := range baseEntityList[4:] {
//			Quad.Remove(e)
//		}
//
//		for _, c := range Quad.children {
//			if len(c.children) != 0 {
//				t.Fail()
//			}
//		}
//	})
//
//}
//
//func TestQuadGo_Retrieve(t *testing.T) {
//	Quad := New(w, h, 4)
//
//	for _, b := range baseBoundsList {
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//	}
//
//	t.Run("Basic Retrieve test", func(t *testing.T) {
//		e := Quad.Retrieve(75, 75)
//		if len(e) == 0 {
//			t.Fail()
//		}
//	})
//
//	//Quad.Insert(400, 80, 550, 280)
//	//Quad.Insert(50, 800, 70, 850)
//	//
//	//t.Run("Retrieve only Bottom left test", func(t *testing.T) {
//	//	e := Quad.Retrieve(0, 0, 50, 50)
//	//	if len(e) > 1 || len(e) == 0 {
//	//		t.Fail()
//	//	}
//	//})
//}
//
//func TestQuadGo_IsEntity(t *testing.T) {
//	Quad := New(w, h, 4)
//
//	for _, e := range baseEntityList {
//		Quad.InsertEntity(e)
//	}
//
//	t.Run("basic isEntity test", func(t *testing.T) {
//		if !Quad.IsEntity(baseEntityList[4]) {
//			t.Fail()
//		}
//	})
//}
//
//func TestQuadGo_IsIntersect(t *testing.T) {
//	Quad := New(w, h, 1)
//
//	for _, b := range baseBoundsList {
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//	}
//
//	t.Run("Basic Intersect test", func(t *testing.T) {
//		if !Quad.IsIntersect(30, 30) {
//			t.Fail()
//		}
//	})
//
//	t.Run("Not Intersected test", func(t *testing.T) {
//		if Quad.IsIntersect(-40, -40) {
//			t.Fail()
//		}
//	})
//
//	t.Run("Intersect topLeft test", func(t *testing.T) {
//		if !Quad.IsIntersect(0, Quad.bounds.Width/2+20, 20, Quad.bounds.Height-20) {
//			t.Fail()
//		}
//	})
//
//	t.Run("Intersect TopRight test", func(t *testing.T) {
//		if !Quad.IsIntersect(Quad.bounds.Width/2+20, Quad.bounds.Height/2+20, Quad.bounds.Width/2+70, Quad.bounds.Height/2+70) {
//			t.Fail()
//		}
//	})
//
//	t.Run("Intersect BottomRight test", func(t *testing.T) {
//		if !Quad.IsIntersect(Quad.bounds.Width/2+20, 0, Quad.bounds.Width/2+70, 50) {
//			t.Fail()
//		}
//	})
//}
//
//func TestQuadGo_Intersects(t *testing.T) {
//	Quad := New(w, h, 1)
//
//	for _, b := range baseBoundsList {
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//	}
//
//	t.Run("Basic Intersects test", func(t *testing.T) {
//		if len(Quad.Intersects(0, 0, 50, 50)) == 0 {
//			t.Fail()
//		}
//	})
//}
//
//func BenchmarkNewQuadGo(b *testing.B) {
//	b.Run("NewQuadGo pass bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			q := New(w, h, 25)
//			if q == nil {
//				b.Fail()
//			}
//		}
//	})
//}
//
//func BenchmarkQuadGo_Insert(b *testing.B) {
//	Quad := New(w, h, 1)
//
//	for n := 0; n < b.N; n++ {
//		for _, b := range baseBoundsList {
//			Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//		}
//	}
//}
//
//func BenchmarkQuadGo_IsIntersect(b *testing.B) {
//	Quad := New(w, h, 25)
//
//	for _, b := range baseBoundsList {
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//	}
//
//	b.Run("IsIntersect true bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			if !Quad.IsIntersect(35, 70, 85, 150) {
//				b.Fail()
//			}
//		}
//	})
//
//	b.Run("IsIntersect false bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			if Quad.IsIntersect(-20, -50, -10, -20) {
//				b.Fail()
//			}
//		}
//	})
//}
//
//func BenchmarkQuadGo_Intersects(b *testing.B) {
//	Quad := New(w, h, 25)
//
//	for _, b := range baseBoundsList {
//		Quad.Insert(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
//	}
//
//	b.Run("Basic intersets bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			if len(Quad.Intersects(25, 50, 50, 75)) == 0 {
//				b.Fail()
//			}
//		}
//	})
//}
//
//func BenchmarkQuadGo_IsEntity(b *testing.B) {
//	Quad := New(w, h, 4)
//
//	for _, e := range baseEntityList {
//		Quad.InsertEntity(e)
//	}
//
//	b.Run("IsEntity Pass bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			if !Quad.IsEntity(baseEntityList[4]) {
//				b.Fail()
//			}
//		}
//	})
//
//	b.Run("IsEntity fail bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			if Quad.IsEntity(&Entity{}) {
//				b.Fail()
//			}
//		}
//	})
//}
//
//func BenchmarkQuadGo_Remove(b *testing.B) {
//	Quad := New(w, h, 5)
//
//	for _, e := range baseEntityList {
//		Quad.InsertEntity(e)
//	}
//
//	b.Run("Remove pass bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			Quad.Remove(baseEntityList[6])
//		}
//	})
//
//	b.Run("Collapse pass bench", func(b *testing.B) {
//		for n := 0; n < b.N; n++ {
//			for _, e := range baseEntityList[4:] {
//				Quad.Remove(e)
//			}
//		}
//	})
//}
