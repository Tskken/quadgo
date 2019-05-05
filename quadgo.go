package QuadGo

import "errors"

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
func NewQuadGO(maxEntities int, rootBounds Bounder) *QuadGo {
	return &QuadGo{
		root: &node{
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

// Retrieve retrieves all entities that are contained in all bounds the given entity fits with in.
func (q *QuadGo) Retrieve(entity Entity) []Entity {
	// retrieve entities for quadtree
	return q.root.retrieve(entity)
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

// node is the container that holds the branch and leaf data for the tree
type node struct {
	bounds   Bounder
	entities []Entity
	children map[quadrant]*node
}

// retrieve finds any entities that are contained in the bounding box the given entity fits in and then returns them
func (n *node) retrieve(entity Entity) (entities []Entity) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// find quadrant the given entity fits in to
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

// insert inserts a given entity in to the quadtree
func (n *node) insert(entity Entity, maxEntities int) {
	// Check if you are on a leaf node
	if len(n.children) > 0 {
		// Find quadrant to insert in to
		if node := n.getQuadrant(entity); node != nil {
			// Insert in to next node
			node.insert(entity, maxEntities)
		} else {
			panic(errors.New("entity does not fit in the max bounds of the tree"))
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
				// Find quadrant to insert entity in to
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

// getQuadrant gets the node for the quadrant the given bounds fits within
func (n *node) getQuadrant(bounds Bounder) *node {
	// get index to quadrant the bounds fits with in
	index := getQuadrant(n.bounds, bounds)
	if index == none {
		return nil
	}

	// return child node for quadrant index
	return n.children[index]
}

// split creates the children for a node by subdividing the nodes boundaries in to 4 even quadrants
func (n *node) split() {
	// new width for child
	subWidth := n.bounds.W() / 2
	// new height for the child
	subHeight := n.bounds.H() / 2
	// nodes bottom left xy coordinates
	x, y := n.bounds.Min().XY()

	// Bottom Left child node
	n.children[bottomLeft] = &node{
		bounds:   NewBounds(x, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Bottom Right child node
	n.children[bottomRight] = &node{
		bounds:   NewBounds(x+subWidth, y, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Top Left child node
	n.children[topLeft] = &node{
		bounds:   NewBounds(x, y+subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}

	// Top Right child node
	n.children[topRight] = &node{
		bounds:   NewBounds(x+subWidth, y+subHeight, subWidth, subHeight),
		entities: make([]Entity, 0),
		children: make(map[quadrant]*node),
	}
}

// getQuadrant returns the quadrant ware the given entityBounds fits with in the given nodeBounds
func getQuadrant(nodeBounds, entityBounds Bounder) quadrant {
	// get the center coordinates for the node bounds
	centerX, centerY := nodeBounds.Center().XY()
	// get the min and max coordinates for the entity bounds
	min, max := entityBounds.Bounds()

	// check if entity fits in the bottom of the node
	bot := min.Y() < centerY && max.Y() <= centerY
	// check if entity fits int he top of node
	top := min.Y() > centerY

	// check if entity fits in the left side of node
	left := min.X() < centerX && max.X() <= centerX
	// check if entity fits in the right side of node
	right := min.X() > centerX

	// return ware the given entity fits in node
	// none means it couldn't fit in node
	switch {
	case bot && left:
		return bottomLeft
	case bot && right:
		return bottomRight
	case top && left:
		return topLeft
	case top && right:
		return topRight
	default:
		return none
	}
}
