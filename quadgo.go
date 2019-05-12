package QuadGo

import "errors"

// constant values for child quadrant map.
const (
	bottomLeft quadrant = iota
	bottomRight
	topLeft
	topRight
)

type quadrant uint8

// QuadGo - Base Quadtree data structure.
type QuadGo struct {
	// Contains all unexported fields
	*node
}

// NewQuadGo creates the basic QuadGo data structure from the given information.
//
// - maxEntities: max number of Entities per node.
//
// - rootBounds: the max bounds of the tree.
func NewQuadGo(maxEntities int, rootBounds Bounds) (*QuadGo, error) {
	if maxEntities <= 0 || rootBounds.validate() {
		return nil, errors.New("given values are not valid")
	}
	return &QuadGo{
		node: &node{
			parent:   nil,
			bounds:   rootBounds,
			entities: make([]Entity, 0, maxEntities),
			children: make([]*node, 0, 4),
		},
	}, nil
}

// Insert inserts an Entity in to the quadtree.
func (q *QuadGo) Insert(bounds Bounds, object interface{}) {
	// insert in to quadtree
	q.insert(Entity{bounds: bounds, object: object})
}

func (q *QuadGo) InsertEntity(entity Entity) {
	q.insert(entity)
}

// Remove removes the given Entity from the quadtree.
func (q *QuadGo) Remove(entity Entity) {
	// remove from quadtree
	q.remove(entity)
}

// Retrieve retrieves all entities that are contained in all bounds the given Entity fits with in.
func (q *QuadGo) Retrieve(bounds Bounds) []Entity {
	// retrieve entities for quadtree
	return q.retrieve(bounds)
}

// IsEntity checks if a given Entity exists in the quadtree.
func (q *QuadGo) IsEntity(entity Entity) bool {
	return q.isEntity(entity)
}

// IsIntersect gets all entities within the bounds that the given Entity fits within and then checks if
// any of the entities intersect with the given Entity.
func (q *QuadGo) IsIntersect(bounds Bounds) bool {
	// check all entities returned from retrieve for if they intersect
	for _, e := range q.Retrieve(bounds) {
		// check for intersect
		if e.bounds.isIntersect(bounds) {
			return true
		}
	}
	return false
}

// Intersects returns a list of all entities the given Entity intersects with.
func (q *QuadGo) Intersects(bounds Bounds) (entities []Entity) {
	// check all entities returned from retrieve for if they intersect
	for _, e := range q.Retrieve(bounds) {
		// add to list if they intersect
		if e.bounds.isIntersect(bounds) {
			entities = append(entities, e)
		}
	}
	return
}

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent   *node
	bounds   Bounds
	entities []Entity
	children []*node
}

// retrieve finds any entities that are contained in the bounding box the given Entity fits in and then returns them.
func (n *node) retrieve(bounds Bounds) (entities []Entity) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// isEntity quadrant the given Entity fits in to
		// - if node is nil returns. Entity could not fit in tree
		if node := n.getQuadrant(bounds); node != nil {
			// add all entities from found quadrant to list
			entities = append(entities, node.retrieve(bounds)...)
		}
	} else {
		// return entities from leaf
		return n.entities
	}
	return
}

// insert inserts a given Entity in to the quadtree.
func (n *node) insert(entity Entity) {
	// Check if you are on a leaf node
	if len(n.children) > 0 {
		// IsEntity quadrant to insert in to
		if node := n.getQuadrant(entity.bounds); node != nil {
			// Insert in to next node
			node.insert(entity)
		}
	} else {
		// Check if a split is needed
		if len(n.entities)+1 > cap(n.entities) {
			// create next leaf nodes
			n.split()

			// loop through all entities to add them to there appropriate child node
			for _, e := range append(n.entities, entity) {
				// IsEntity quadrant to insert Entity in to
				// Nil means it didn't fit in to any quadrant
				if node := n.getQuadrant(e.bounds); node != nil {
					// insert Entity to new child
					node.insert(e)
				}
			}
			// clear entities for branch node
			n.entities = make([]Entity, 0, cap(n.entities))
		} else {
			// Add Entity to node
			n.entities = append(n.entities, entity)
		}
	}
}

// remove removes the given Entity from the quadtree.
func (n *node) remove(entity Entity) {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// not on a leaf, get next quadrant
		if node := n.getQuadrant(entity.bounds); node != nil {
			node.remove(entity)
		}
	} else {
		// check entities in leaf for given Entity
		for i, e := range n.entities {
			// check if given Entity is the same as node Entity
			if e == entity {
				// check if removal would make the leaf have no entities
				if len(n.entities) == 1 {
					// set node entities to nil
					n.entities = make([]Entity, 0, cap(n.entities))

					// check if children can be collapsed in to parent node
					n.parent.collapse()
				} else {
					// remove Entity from node
					n.entities = append(n.entities[:i], n.entities[i+1:]...)
				}
			}
		}
	}
}

// collapse checks if a parent's children hold less entities then the set maxEntities count.
// if the count is less then maxEntities it collapses all children in to the parent node, copying
// all of there entities to the parent node and setting the children to nil.
func (n *node) collapse() {
	// create base counter for children Entity count
	eCount := 0
	for _, c := range n.children {
		// add children's Entity count to counter
		eCount += len(c.entities)
	}

	// check if the total number of entities in the nodes children is
	// less then the max number of entities allowed in an node
	if eCount < cap(n.entities) {
		// move children entities to parent node
		for _, c := range n.children {
			n.entities = append(n.entities, c.entities...)
		}

		// reset children
		n.children = make([]*node, 0, 4)
	}
}

// isEntity returns if a given Entity exists in the quadtree.
func (n *node) isEntity(entity Entity) bool {
	// find all entities that could match given Entity
	for _, e := range n.retrieve(entity.bounds) {
		// check if given Entity equals Entity
		if e.bounds == entity.bounds {
			return true
		}
	}

	return false
}

// split creates the children for a node by subdividing the nodes boundaries in to 4 even quadrants.
func (n *node) split() {
	center := n.bounds.center()

	// Bottom Left child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(n.bounds.min.x, n.bounds.min.y, center.x, center.y),
		entities: make([]Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Bottom Right child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(center.x, n.bounds.min.y, n.bounds.max.x, center.y),
		entities: make([]Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Top Left child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(n.bounds.min.x, center.y, center.x, n.bounds.max.y),
		entities: make([]Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Top Right child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(center.x, center.y, n.bounds.max.x, n.bounds.max.y),
		entities: make([]Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})
}

// getQuadrant returns the quadrant ware the given entityBounds fits with in the given nodeBounds.
func (n *node) getQuadrant(bounds Bounds) *node {
	// get the center coordinates for the node bounds
	center := n.bounds.center()

	switch {
	case (bounds.min.x < center.x && bounds.max.x <= center.x) && (bounds.min.y < center.y && bounds.max.y <= center.y):
		return n.children[bottomLeft]
	case (bounds.min.x >= center.x) && (bounds.min.y < center.y && bounds.max.y <= center.y):
		return n.children[bottomRight]
	case (bounds.min.x < center.x && bounds.max.x <= center.x) && (bounds.min.y >= center.y):
		return n.children[topLeft]
	case (bounds.min.x >= center.x) && (bounds.min.y >= center.y):
		return n.children[topRight]
	default:
		return nil
	}
}
