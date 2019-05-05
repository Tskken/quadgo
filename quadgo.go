package QuadGo

type QuadGo struct {
	root        *node
	maxEntities int
	size        int
}

func NewQuadGO(maxEntities int, rootBounds Bounder) *QuadGo {
	return &QuadGo{
		root:&node{
			bounds:   rootBounds,
			entities: make([]Entity, 0),
			children: make([]*node, 0, 4),
			level:    1,
		},
		maxEntities: maxEntities,
		size:        1,
	}
}

func (q *QuadGo) Insert(entity Entity) {
	q.size++
	q.root.insert(entity, q.maxEntities)
}

func (q *QuadGo) Retrieve(entity Entity) []Entity {
	return q.root.retrieve(entity)
}

func (q *QuadGo) IsIntersect(entity Entity) bool {
	for _, e := range q.Retrieve(entity) {
		if e.IsIntersect(entity) {
			return true
		}
	}
	return false
}

func (q *QuadGo) Intersects(entity Entity) (entities []Entity) {
	for _, e := range q.Retrieve(entity) {
		if e.IsIntersect(entity) {
			entities = append(entities, e)
		}
	}
	return
}

type node struct {
	bounds   Bounder
	entities []Entity
	children []*node
	level    int
}

func (n *node) retrieve(entity Entity) (entities []Entity) {
	if len(n.children) >= 0 {
		if node := n.getQuadrant(entity); node != nil {
			entities = append(entities, n.retrieve(entity)...)
		}
	} else {
		return n.entities
	}
	return
}

func (n *node) insert(entity Entity, maxEntities int) {
	// Check if you are on a leaf or branch
	if len(n.children) > 0 {
		// Find quadrant to insert in to
		if node := n.getQuadrant(entity); node != nil {
			// Insert in to next node
			node.insert(entity, maxEntities)
		}
	} else {
		// Add entity to node
		n.entities = append(n.entities, entity)

		if len(n.entities) > maxEntities {
			n.split()
			for i := 0; i < len(n.entities);{
				// Find quadrant to insert in to
				// Nil means it didn't fit in to any quadrant
				if node := n.getQuadrant(n.entities[i]); node != nil {
					splice := n.entities[i]
					n.entities = append(n.entities[:i], n.entities[i+1:]...)

					node.insert(splice, maxEntities)
				} else {
					i++
				}
			}
		}
	}
}

func (n *node) getQuadrant(bounds Bounder) *node {
	index := getQuadrant(n.bounds, bounds)
	if index == -1 {
		return nil
	}
	return n.children[index]
}

func (n *node) split() {
	nextLevel := n.level +1
	subWidth := n.bounds.W() / 2
	subHeight := n.bounds.H() / 2
	x, y := n.bounds.Min().XY()

	n.children = append(n.children, &node{
		bounds:   NewBounds(x, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make([]*node, 0 , 4),
		level:    nextLevel,
	})

	n.children = append(n.children, &node{
		bounds:   NewBounds(x + subWidth, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make([]*node, 0, 4),
		level:    nextLevel,
	})

	n.children = append(n.children, &node{
		bounds:   NewBounds(x, y + subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make([]*node, 0, 4),
		level:    nextLevel,
	})

	n.children = append(n.children, &node{
		bounds:   NewBounds(x+subWidth, y+subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make([]*node, 0, 4),
		level:    nextLevel,
	})
}

func getQuadrant(nodeBounds, entityBounds Bounder) int {
	centerX, centerY := nodeBounds.Center().XY()
	min, max := entityBounds.Bounds()

	bot := min.Y() < centerY && max.Y() < centerY
	top := min.Y() > centerY
	left := min.X() < centerX && max.X() < centerX
	right := min.Y() > centerY

	switch {
	case bot && left:
		return 0
	case bot && right:
		return 1
	case top && left:
		return 2
	case top && right:
		return 3
	default:
		return -1
	}
}

