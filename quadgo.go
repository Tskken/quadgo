package quadgo

import (
	"errors"
	"math"
)

// constant values for child quadrants.
const (
	bottomLeft quadrant = iota
	bottomRight
	topLeft
	topRight
)

type quadrant uint8

// QuadGo - Base Quadtree data structure.
type QuadGo struct {
	*node
}

// NewQuadGo creates the basic QuadGo data structure from the given information.
//
// - maxEntities: max number of Entities per node.
//
// - screenWidth: Width of the screen or map that will encompass all Bounds and objects.
//
// - screenHeight: Height of the screen or map that will encompass all Bounds and objects.
func NewQuadGo(maxEntities int, screenWidth, screenHeight float64) (*QuadGo, error) {
	if maxEntities <= 0 {
		return nil, errors.New("given values are not valid")
	}
	return &QuadGo{
		node: &node{
			parent:   nil,
			bounds:   NewBounds(0, 0, math.Abs(screenWidth), math.Abs(screenHeight)),
			entities: make([]*Entity, 0, maxEntities),
			children: make([]*node, 0, 4),
		},
	}, nil
}

// Insert inserts a Bounds in to the quadtree with a corresponding Object.
//
// The Object is any data type you may want to store in the quadtree that is not a Bounds.
// When searching the tree it will return a Entity which holds the given Bounds and the Object provide.
//
// If you do not want to add an Object to the tree you can just put nil.
func (q *QuadGo) Insert(bounds Bounds, object ...interface{}) {
	// insert in to quadtree
	q.insert(&Entity{Bounds: bounds, Object: object})
}

// InsertEntity inserts an entity in to the quadtree.
func (q *QuadGo) InsertEntity(entity *Entity) {
	q.insert(entity)
}

// Remove removes the given Entity from the quadtree.
func (q *QuadGo) Remove(entity *Entity) {
	// remove from quadtree
	q.remove(entity)
}

// Retrieve returns a list of all entities that are with in a nodes Bounds that the given Bounds fits with in.
func (q *QuadGo) Retrieve(bounds Bounds) []*Entity {
	// retrieve entities for quadtree
	return q.retrieve(bounds)
}

// IsEntity checks if a given Entity exists in the quadtree.
//
// Note: This function currently is very slow for unknown reasons and it is advised to just not use it.
// If you are going to use IsEntity() for something understand it may slow down performance significantly.
// In a future update (likely V2.0.1) I hope to fix this issue but for know be noted this function is not advised to be
// used.
func (q *QuadGo) IsEntity(entity *Entity) bool {
	return q.isEntity(entity)
}

// IsIntersect takes a Bounds and returns if it intersect with any entity in the quadtree.
func (q *QuadGo) IsIntersect(bounds Bounds) bool {
	entities := q.retrieve(bounds)
	// check all entities returned from retrieve for if they intersect
	for i := range entities {
		// check for intersect
		if entities[i].Bounds.isIntersect(bounds) {
			return true
		}
	}
	return false
}

// Intersects takes a Bounds and returns a list of all entities it intersects with.
func (q *QuadGo) Intersects(bounds Bounds) (intersects []*Entity) {
	entities := q.retrieve(bounds)
	// check all entities returned from retrieve for if they intersect
	for i := range entities {
		// add to list if they intersect
		if entities[i].Bounds.isIntersect(bounds) {
			intersects = append(intersects, entities[i])
		}
	}
	return
}

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent   *node
	bounds   Bounds
	entities []*Entity
	children []*node
}

// retrieve finds all of the entities with in a the nodes Bounds that the given Bounds can fit with in.
func (n *node) retrieve(bounds Bounds) []*Entity {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// isEntity quadrant the given Entity fits in to
		// - if node is nil returns. Entity could not fit in tree
		if node := n.getQuadrant(bounds); node != nil {
			// add all entities from found quadrant to list
			return node.retrieve(bounds)
		}
	} else {
		// return entities from leaf
		return n.entities
	}
	return nil
}

// insert inserts a given Entity in to the quadtree.
func (n *node) insert(entity *Entity) {
	// Check if you are on a leaf node
	if len(n.children) > 0 {
		// IsEntity quadrant to insert in to
		if node := n.getQuadrant(entity.Bounds); node != nil {
			// Insert in to next node
			node.insert(entity)
		}
	} else {
		// Check if a split is needed
		if len(n.entities)+1 > cap(n.entities) {
			// create next leaf nodes
			n.split()

			entities := append(n.entities, entity)

			// loop through all entities to add them to there appropriate child node
			for i := range entities {
				// IsEntity quadrant to insert Entity in to
				// Nil means it didn't fit in to any quadrant
				if node := n.getQuadrant(entities[i].Bounds); node != nil {
					// insert Entity to new child
					node.insert(entities[i])
				}
			}
			// clear entities for branch node
			n.entities = make([]*Entity, 0, cap(n.entities))
		} else {
			// Add Entity to node
			n.entities = append(n.entities, entity)
		}
	}
}

// remove removes the given Entity from the quadtree.
func (n *node) remove(entity *Entity) {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// not on a leaf, get next quadrant
		if node := n.getQuadrant(entity.Bounds); node != nil {
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
					n.entities = make([]*Entity, 0, cap(n.entities))

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
	for i := range n.children {
		// add children's Entity count to counter
		eCount += len(n.children[i].entities)
	}

	// check if the total number of entities in the nodes children is
	// less then the max number of entities allowed in an node
	if eCount < cap(n.entities) {
		// move children entities to parent node
		for i := range n.children {
			n.entities = append(n.entities, n.children[i].entities...)
		}

		// reset children
		n.children = make([]*node, 0, 4)
	}
}

// isEntity returns if a given Entity exists in the quadtree.
func (n *node) isEntity(entity *Entity) bool {
	entities := n.retrieve(entity.Bounds)
	// find all entities that could match given Entity
	for i := range entities {
		// check if given Entity equals Entity
		if entities[i] == entity {
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
		entities: make([]*Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Bottom Right child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(center.x, n.bounds.min.y, n.bounds.max.x, center.y),
		entities: make([]*Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Top Left child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(n.bounds.min.x, center.y, center.x, n.bounds.max.y),
		entities: make([]*Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})

	// Top Right child node
	n.children = append(n.children, &node{
		parent:   n,
		bounds:   NewBounds(center.x, center.y, n.bounds.max.x, n.bounds.max.y),
		entities: make([]*Entity, 0, cap(n.entities)),
		children: make([]*node, 0, 4),
	})
}

// getQuadrant returns the nodes child node that the given Bounds fits with in.
func (n *node) getQuadrant(bounds Bounds) *node {
	// get the center coordinates for the node Bounds
	center := n.bounds.center()

	if (bounds.min.x < center.x && bounds.max.x <= center.x) && (bounds.min.y < center.y && bounds.max.y <= center.y) {
		return n.children[bottomLeft]
	} else if (bounds.min.x >= center.x) && (bounds.min.y < center.y && bounds.max.y <= center.y){
		return n.children[bottomRight]
	} else if (bounds.min.x < center.x && bounds.max.x <= center.x) && (bounds.min.y >= center.y){
		return n.children[topLeft]
	} else if (bounds.min.x >= center.x) && (bounds.min.y >= center.y){
		return n.children[topRight]
	}

	return nil
}
