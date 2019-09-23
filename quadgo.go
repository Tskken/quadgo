package quadgo

import (
	"errors"
	"fmt"
)

// quadrant type for iota child quadrants.
type quadrant uint8

// constant values for child quadrants.
const (
	topLeft quadrant = iota
	topRight
	bottomLeft
	bottomRight
)

// Option function type for setting the Options of a new tree.
type Option func(*Options)

// Options struct which holds all the information for creating a new quadtree with its given information.
type Options struct {
	MaxEntities uint64
	MaxDepth    uint16
}

// defaultOptions for QuadGo
var defaultOption = &Options{
	MaxEntities: 10,
	MaxDepth:    2,
}

// SetMaxEntities sets the max number of entities per each node in the new tree.
func SetMaxEntities(maxEntities uint64) Option {
	return func(o *Options) {
		o.MaxEntities = maxEntities
	}
}

// SetMaxDepth sets the max depth that the tree can splitAndMove to.
func SetMaxDepth(maxDepth uint16) Option {
	return func(o *Options) {
		o.MaxDepth = maxDepth
	}
}

// QuadGo - Base Quadtree data structure.
type QuadGo struct {
	*node

	maxDepth uint16
}

// New creates the basic QuadGo instance.
//
// New requires a width and a height but can also be given any number of other supported option functions.
//
// Example:
//	basic - quadgo.New(800, 600)
// 	with option - quadgo.New(800, 600, SetMaxDepth(5))
//
// QuadGo sets the New defaults for max depth to 2 and max entities to 10.
func New(width, height float64, ops ...Option) *QuadGo {
	// copy defaults
	o := defaultOption

	// update for any given options
	for _, op := range ops {
		op(o)
	}

	// Return new QuadGo instance
	return &QuadGo{
		node: &node{
			parent:   nil,
			bound:    NewBound(0, 0, width, height),
			entities: make(Entities, 0, o.MaxEntities),
			children: make(nodes, 0, 4),
			depth:    0,
		},
		maxDepth: o.MaxDepth,
	}
}

// Insert takes a new entities Min and Max xy bounds and inserts it in to the quadtree.
func (q *QuadGo) Insert(minX, minY, maxX, maxY float64) error {
	return q.insert(NewEntity(minX, minY, maxX, maxY), q.maxDepth)
}

// InsertWithAction takes a new entities Min and Max xy bounds and a Action function and inserts it in to the quadtree.
func (q *QuadGo) InsertWithAction(minX, minY, maxX, maxY float64, action Action) error {
	return q.insert(NewEntityWithAction(minX, minY, maxX, maxY, action), q.maxDepth)
}

// InsertEntities inserts any number of entities in to the quadtree.
//
// This can be used as a second option over Insert if you want to create your entities before adding it to the quadtree,
// or if you need to reenter a entity after removing it from the tree.
func (q *QuadGo) InsertEntities(entities ...*Entity) error {
	// check for no entities given on function call
	if len(entities) == 0 {
		return errors.New("no entities given to QuadGo.InsertEntities()")
	}

	// insert each given entities to the tree
	for _, e := range entities {
		err := q.insert(e, q.maxDepth)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove removes the given Entity from the quadtree.
//
// The given entity only has to have the same data of the entity you want to remove. It does not have to be the exact
// reference to the node you wish to delete.
func (q *QuadGo) Remove(entity *Entity) error {
	return q.remove(entity)
}

// RetrieveFromPoint returns a list of entities that are stored in the node that the given point can be contained within.
//
// If there was no entities in the node for the given point or there was no quadrant for that point it will return an empty slice of entities.
func (q *QuadGo) RetrieveFromPoint(point Point) Entities {
	return q.retrieve(point)
}

// RetrieveFromBound returns a list of entities that are stored in a node that the given bound's center point can be contained within.
//
// If there was no entities in the node for the given bound or there was no quadrant for that bound it will return an empty slice of entities.
func (q *QuadGo) RetrieveFromBound(bound Bound) Entities {
	return q.retrieve(bound.Center)
}

// IsEntity checks if a given entity exists within the tree.
func (q *QuadGo) IsEntity(entity *Entity) bool {
	return q.isEntity(entity)
}

// IsIntersectPoint takes a point and returns if that point intersects any entity within the tree.
func (q *QuadGo) IsIntersectPoint(point Point) bool {
	// get possible entities that the given point could intersect with
	entities := q.retrieve(point)

	// check if any entities returned intersect the given point
	for i := range entities {
		// check for intersect
		if entities[i].IsIntersectPoint(point) {
			return true
		}
	}
	return false
}

// IsIntersectBound take a bound and returns if that bound intersects any entity within the tree.
func (q *QuadGo) IsIntersectBound(bound Bound) bool {
	// get possible entities that the given bound could intersect with
	entities := q.retrieve(bound.Center)

	// check if any entities returned intersect the given bound
	for i := range entities {
		// check for intersect
		if entities[i].IsIntersectBound(bound) {
			return true
		}
	}
	return false
}

// IntersectsPoint takes a point and returns all entities that that point intersects with within the tree.
func (q *QuadGo) IntersectsPoint(point Point) (intersects Entities) {
	// get possible entities the given point could intersect with
	entities := q.retrieve(point)

	// check if any entities returned intersect the given point and if they do add them to the return list
	for i := range entities {
		// add to list if they intersect
		if entities[i].IsIntersectPoint(point) {
			intersects = append(intersects, entities[i])
		}
	}
	return
}

// IntersectsBound takes a bound and returns all entities that that bound intersects with within the tree.
func (q *QuadGo) IntersectsBound(bound Bound) (intersects Entities) {
	// get possible entities the given bound could intersect with
	entities := q.retrieve(bound.Center)

	// check if any entities returned intersect the given bound and if they do add them to the return list
	for i := range entities {
		// add to list if they intersect
		if entities[i].IsIntersectBound(bound) {
			intersects = append(intersects, entities[i])
		}
	}
	return
}

// list of node
type nodes []*node

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent     *node
	bound      Bound
	entities   Entities
	children   nodes
	depth, max uint16
}

// new creates a new node instance for a given bounds taking the parent nodes root.
func (n *node) new(bound Bound) *node {
	return &node{
		parent:   n,
		bound:    bound,
		entities: make(Entities, 0, cap(n.entities)),
		children: make(nodes, 0, 4),
		depth:    n.depth + 1,
	}
}

// retrieve finds all of the entities with in a quadrant that the given point fits in.
func (n *node) retrieve(point Point) (e Entities) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// get the quadrant that the point fits in and go to that next node
		if node := n.getQuadrant(point); node != nil {
			return node.retrieve(point)
		}
		// return an empty list for no quadrant found for given point
		return
	}

	// return entities from leaf
	return n.entities
}

// insert inserts a given entity in to the quadtree.
func (n *node) insert(entity *Entity, maxDepth uint16) error {
	// check if you are on a leaf node or at max depth of the tree
	if len(n.children) > 0 && n.depth <= n.max {
		// get the next node that the given entity fits in and attempt to insert it
		if node := n.getQuadrant(entity.Center); node != nil {
			return node.insert(entity, maxDepth)
		}

		// return an error for no quadrants found
		return errors.New("returned node from getQuadrant was nil")
	}

	// check if a splitAndMove is needed
	if len(n.entities)+1 > cap(n.entities) && n.depth < n.max {
		// split node in to child nodes
		n.split()

		// move this nodes entities to the children nodes
		return n.MoveEntities(append(n.entities, entity), maxDepth)
	}

	// add Entity to node
	n.entities = append(n.entities, entity)
	return nil
}

// remove removes the given Entity from the quadtree.
func (n *node) remove(entity *Entity) error {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// get the next node that the given entity fits in and attempt to remove it
		if node := n.getQuadrant(entity.Center); node != nil {
			return node.remove(entity)
		}
		// return an error for no quadrants found for given entity
		return fmt.Errorf("could not find a quadrent for the given entity: %v", entity)
	}

	entities, err := n.entities.Remove(entity)
	if err != nil {
		return err
	}

	n.entities = entities

	if n.parent != nil {
		// check if children can be collapsed in to parent node
		if n.parent.shouldCollapse() {
			n.parent.collapse()
		}
	}

	return nil
}

// shouldCollapse checks if the nodes children should be collapsed in to the parent.
func (n *node) shouldCollapse() bool {
	// create base counter for children entity count
	eCount := 0

	// count up total entities in children
	for i := range n.children {
		eCount += len(n.children[i].entities)
	}

	return eCount < cap(n.entities)
}

// collapse takes all entities from the children nodes and moves them to the parent and then removes the children.
func (n *node) collapse() {
	// move children entities to parent node
	for i := range n.children {
		n.entities = append(n.entities, n.children[i].entities...)
	}

	// remove children
	n.children = n.children[:0]
}

// isEntity returns if a given entity exists in the tree.
func (n *node) isEntity(entity *Entity) bool {
	// get entities from a node that the entity.Center can fit in
	entities := n.retrieve(entity.Center)

	// check each entity for if it is equal to given entity
	for i := range entities {
		// check if given Entity equals given entity
		if entities[i].IsEqual(entity) {
			return true
		}
	}

	return false
}

// split creates the children node for this node.
func (n *node) split() {
	// Top Left child node
	n.children = append(n.children, n.new(NewBound(n.bound.Min.X, n.bound.Center.Y, n.bound.Center.X, n.bound.Max.Y)))

	// Top Right child node
	n.children = append(n.children, n.new(NewBound(n.bound.Center.X, n.bound.Center.Y, n.bound.Max.X, n.bound.Max.Y)))

	// Bottom Left child node
	n.children = append(n.children, n.new(NewBound(n.bound.Min.X, n.bound.Min.Y, n.bound.Center.X, n.bound.Center.Y)))

	// Bottom Right child node
	n.children = append(n.children, n.new(NewBound(n.bound.Center.X, n.bound.Min.Y, n.bound.Max.X, n.bound.Center.Y)))
}

func (n *node) MoveEntities(entities Entities, maxDepth uint16) error {
	// loop through all entities to add them to there appropriate child node
	for i := range entities {
		// get the next node that the given entity fits in and insert it
		err := n.insert(entities[i], maxDepth)
		if err != nil {
			return err
		}
	}

	// clear entities for branch node
	n.entities = n.entities[:0]
	return nil
}

// getQuadrant returns this nodes child node that the given point fits within.
func (n *node) getQuadrant(point Point) *node {
	if n.bound.Center.X < point.X && n.bound.Center.Y >= point.Y {
		return n.children[topLeft] // return top left
	} else if n.bound.Center.X >= point.X && n.bound.Center.Y >= point.Y {
		return n.children[topRight] // return top right
	} else if n.bound.Center.X < point.X && n.bound.Center.Y < point.Y {
		return n.children[bottomLeft] // return bottom left
	} else {
		return n.children[bottomRight] // return bottom right
	}
}
