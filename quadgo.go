package QuadGo

// constant values for child quadrant map.
const (
	bottomLeft quadrant = iota
	bottomRight
	topLeft
	topRight
	none
)

type quadrant byte

// QuadGo - Base Quadtree data structure.
type QuadGo struct {
	// Contains all unexported fields
	root        *node
	maxEntities int
}

// NewQuadGo creates the basic QuadGo data structure from the given information.
//
// - maxEntities: max number of Entities per node.
//
// - rootBounds: the max bounds of the tree.
func NewQuadGo(maxEntities int, rootBounds Bounder) *QuadGo {
	return &QuadGo{
		root: &node{
			parent:nil,
			bounds:   rootBounds,
			entities: make([]Entity, 0),
			children: make(map[quadrant]*node),
		},
		maxEntities: maxEntities,
	}
}

// Insert inserts an Entity in to the quadtree.
func (q *QuadGo) Insert(entity Entity) {
	// insert in to quadtree
	q.root.insert(entity, q.maxEntities)
}

// Remove removes the given entity from the quadtree.
func (q *QuadGo) Remove(entity Entity) {
	// remove from quadtree
	q.root.remove(entity, q.maxEntities)
}

// Retrieve retrieves all entities that are contained in all bounds the given entity fits with in.
func (q *QuadGo) Retrieve(entity Entity) []Entity {
	// retrieve entities for quadtree
	return q.root.retrieve(entity)
}

// IsEntity checks if a given entity exists in the quadtree.
func (q *QuadGo) IsEntity(entity Entity) bool {
	return q.root.isEntity(entity)
}

// IsIntersect gets all entities within the bounds that the given entity fits within and then checks if
// any of the entities intersect with the given entity.
func (q *QuadGo) IsIntersect(entity Entity) bool {
	// check all entities returned from retrieve for if they intersect
	for _, e := range q.Retrieve(entity) {
		// check for intersect
		if e.IsIntersect(entity) {
			return true
		}
	}
	return false
}

// Intersects returns a list of all entities the given entity intersects with.
func (q *QuadGo) Intersects(entity Entity) (entities []Entity) {
	// check all entities returned from retrieve for if they intersect
	for _, e := range q.Retrieve(entity) {
		// add to list if they intersect
		if e.IsIntersect(entity) {
			entities = append(entities, e)
		}
	}
	return
}

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent *node
	bounds   Bounder
	entities []Entity
	children map[quadrant]*node
}

// retrieve finds any entities that are contained in the bounding box the given entity fits in and then returns them.
func (n *node) retrieve(entity Entity) (entities []Entity) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// isEntity quadrant the given entity fits in to
		// - if node is nil returns. Entity could not fit in tree
		if node := n.getQuadrant(entity); node != nil {
			// add all entities from found quadrant to list
			entities = append(entities, node.retrieve(entity)...)
		}
	} else {
		return n.entities
	}
	return
}

// insert inserts a given entity in to the quadtree.
func (n *node) insert(entity Entity, maxEntities int) {
	// Check if you are on a leaf node
	if len(n.children) > 0 {
		// IsEntity quadrant to insert in to
		if node := n.getQuadrant(entity); node != nil {
			// Insert in to next node
			node.insert(entity, maxEntities)
		}
	} else {
		// Add entity to node
		n.entities = append(n.entities, entity)

		// Check if a split is needed
		if len(n.entities) > maxEntities {
			// create next leaf nodes
			n.split()

			// loop through all entities to add them to there appropriate child node
			for _, e := range n.entities {
				// IsEntity quadrant to insert entity in to
				// Nil means it didn't fit in to any quadrant
				if node := n.getQuadrant(e); node != nil {
					// insert entity to new child
					node.insert(e, maxEntities)
				}
			}
			// clear entities for branch node
			n.entities = []Entity{}
		}
	}
}

// remove removes the given entity from the quadtree.
func (n *node) remove(entity Entity, maxEntities int) {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// not on a leaf, get next quadrant
		if node := n.getQuadrant(entity); node != nil {
			node.remove(entity, maxEntities)
		}
	} else {
		// check entities in leaf for given entity
		for i, e := range n.entities {
			// check if given entity is the same as node entity
			if e == entity {
				// check if removal would make the leaf have no entities
				if len(n.entities) == 1 {
					// set node entities to nil
					n.entities = make([]Entity, 0)

					n.parent.collapse(maxEntities)
				} else {
					n.entities = append(n.entities[:i], n.entities[i+1:]...)
				}
			}
		}
	}
}

// collapse checks if a parent's children hold less entities then the set maxEntities count.
// if the count is less then maxEntities it collapses all children in to the parent node, copying
// all of there entities to the parent node and setting the children to nil.
func (n *node) collapse(maxEntities int) {
	eCount := 0
	for _, c := range n.children {
		eCount += len(c.entities)
	}

	if eCount < maxEntities {
		for _, c := range n.children {
			n.entities = append(n.entities, c.entities...)
		}

		n.children = make(map[quadrant]*node)
	}
}

// isEntity returns if a given entity exists in the quadtree.
func (n *node) isEntity(entity Entity) bool {
	for _, e := range n.retrieve(entity) {
		if e == entity {
			return true
		}
	}

	return false
}

// getQuadrant gets the node for the quadrant the given bounds fits within.
func (n *node) getQuadrant(bounds Bounder) *node {
	// get index to quadrant the bounds fits with in
	if index := getQuadrant(n.bounds, bounds); index != none {
		// return child node for quadrant index
		return n.children[index]
	}
	return nil
}

// split creates the children for a node by subdividing the nodes boundaries in to 4 even quadrants.
func (n *node) split() {
	// new width for child
	subWidth := n.bounds.W() / 2
	// new height for the child
	subHeight := n.bounds.H() / 2
	// nodes bottom left xy coordinates
	x, y := n.bounds.Min()

	// Bottom Left child node
	n.children[bottomLeft] = &node{
		parent:n,
		bounds:   NewBounds(x, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Bottom Right child node
	n.children[bottomRight] = &node{
		parent:n,
		bounds:   NewBounds(x+subWidth, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Top Left child node
	n.children[topLeft] = &node{
		parent:n,
		bounds:   NewBounds(x, y+subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Top Right child node
	n.children[topRight] = &node{
		parent:n,
		bounds:   NewBounds(x+subWidth, y+subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}
}

// getQuadrant returns the quadrant ware the given entityBounds fits with in the given nodeBounds.
func getQuadrant(nodeBounds, entityBounds Bounder) quadrant {
	// get the center coordinates for the node bounds
	centerX, centerY := nodeBounds.Center()
	// get the min and max coordinates for the entity bounds
	minX, minY, maxX, maxY := entityBounds.Bounds()

	// return ware the given entity fits in node
	// none means it couldn't fit in node
	switch {
	case  (minY < centerY && maxY <= centerY) && (minX < centerX && maxX <= centerX):
		return bottomLeft
	case (minY < centerY && maxY <= centerY) && (minX > centerX):
		return bottomRight
	case (minY > centerY) && (minX < centerX && maxX <= centerX):
		return topLeft
	case (minY > centerY) && (minX > centerX):
		return topRight
	default:
		return none
	}
}
