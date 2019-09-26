package quadgo

import (
	"errors"
	"sync"
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
var defaultOption = Options{
	MaxEntities: 10,
	MaxDepth:    4,
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

	sync.RWMutex
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
		op(&o)
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
func (q *QuadGo) Insert(minX, minY, maxX, maxY float64) {
	q.Lock()
	defer q.Unlock()

	q.insert(NewEntity(minX, minY, maxX, maxY), q.maxDepth)
}

// InsertWithAction takes a new entities Min and Max xy bounds and a Action function and inserts it in to the quadtree.
func (q *QuadGo) InsertWithAction(minX, minY, maxX, maxY float64, action Action) {
	q.Lock()
	defer q.Unlock()

	q.insert(NewEntityWithAction(minX, minY, maxX, maxY, action), q.maxDepth)
}

// InsertEntities inserts any number of entities in to the quadtree.
//
// This can be used as a second option over Insert if you want to create your entities before adding it to the quadtree,
// or if you need to reenter a entity after removing it from the tree.
func (q *QuadGo) InsertEntities(entities ...*Entity) error {
	q.Lock()
	defer q.Unlock()

	// check for no entities given on function call
	if len(entities) == 0 {
		return errors.New("no entities given to QuadGo.InsertEntities()")
	}

	// insert each given entities to the tree
	for _, e := range entities {
		q.insert(e, q.maxDepth)
	}
	return nil
}

// Remove removes the given Entity from the quadtree.
//
// the given entity has to have the same bound size to match. The action function does not need to be the same.
func (q *QuadGo) Remove(entity *Entity) error {
	q.Lock()
	defer q.Unlock()

	return q.remove(entity)
}

// RetrieveFromPoint returns a list of entities that are stored in the node that the given point can be contained within.
//
// If there was no entities in the node for the given point or there was no quadrant for that point it will return an empty slice of entities.
func (q *QuadGo) RetrieveFromPoint(point Point) <-chan Entities {
	out := make(chan Entities)

	go func() {
		q.RLock()
		q.retrieve(point, out)
		q.RUnlock()
	}()

	return out
}

// RetrieveFromBound returns a list of entities that are stored in a node that the given bound's center point can be contained within.
//
// If there was no entities in the node for the given bound or there was no quadrant for that bound it will return an empty slice of entities.
func (q *QuadGo) RetrieveFromBound(bound Bound) <-chan Entities {
	out := make(chan Entities)

	go func() {
		q.RLock()
		q.retrieve(bound.Center, out)
		q.RUnlock()
	}()

	return out
}

// IsEntity checks if a given entity exists within the tree.
func (q *QuadGo) IsEntity(entity *Entity) <-chan bool {
	out := make(chan bool)

	go func() {
		q.RLock()
		q.isEntity(entity, out)
		q.RUnlock()
	}()

	return out
}

// IsIntersectPoint takes a point and returns if that point intersects any entity within the tree.
func (q *QuadGo) IsIntersectPoint(point Point) <-chan bool {
	out := make(chan bool)

	go func() {
		q.RLock()

		// get possible entities that the given point could intersect with
		data := q.RetrieveFromPoint(point)

		entities := <-data

		entities.isIntersectPoint(point, out)

		q.RUnlock()
	}()

	return out
}

// IsIntersectBound take a bound and returns if that bound intersects any entity within the tree.
func (q *QuadGo) IsIntersectBound(bound Bound) <-chan bool {
	out := make(chan bool)

	go func() {
		q.RLock()

		data := q.RetrieveFromBound(bound)

		entities := <-data

		entities.isIntersectBound(bound, out)

		q.RUnlock()
	}()

	return out
}

// IntersectsPoint takes a point and returns all entities that that point intersects with within the tree.
func (q *QuadGo) IntersectsPoint(point Point) <-chan Entities {
	out := make(chan Entities)

	go func() {
		q.RLock()

		data := q.RetrieveFromPoint(point)

		entities := <-data

		entities.intersectsPoint(point, out)

		q.RUnlock()
	}()

	return out
}

// IntersectsBound takes a bound and returns all entities that that bound intersects with within the tree.
func (q *QuadGo) IntersectsBound(bound Bound) <-chan Entities {
	out := make(chan Entities)

	go func() {
		q.RLock()

		data := q.RetrieveFromBound(bound)

		entities := <-data

		entities.intersectsBound(bound, out)

		q.RUnlock()
	}()

	return out
}

// list of node
type nodes []*node

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent   *node
	bound    Bound
	entities Entities
	children nodes
	depth    uint16
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
func (n *node) retrieve(point Point, out chan Entities) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// get the quadrant that the point fits in and go to that next node
		n.getQuadrant(point).retrieve(point, out)
	}

	// return entities from leaf
	out <- n.entities
}

// insert inserts a given entity in to the quadtree.
func (n *node) insert(entity *Entity, maxDepth uint16) {
	// check if you are on a leaf node or at max depth of the tree
	if len(n.children) > 0 {
		// get the next node that the given entity fits in and attempt to insert it
		n.getQuadrant(entity.Center).insert(entity, maxDepth)
		return
	}

	// check if a splitAndMove is needed
	if len(n.entities)+1 > cap(n.entities) && n.depth < maxDepth {
		// split node in to child nodes
		n.split()

		// move this nodes entities to the children nodes
		n.MoveEntities(append(n.entities, entity), maxDepth)
		return
	}

	// add Entity to node
	n.entities = append(n.entities, entity)
}

// remove removes the given Entity from the quadtree.
func (n *node) remove(entity *Entity) error {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// get the next node that the given entity fits in and attempt to remove it
		return n.getQuadrant(entity.Center).remove(entity)
	}

	entities, err := n.entities.FindAndRemove(entity)
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

	return eCount <= cap(n.entities)
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
func (n *node) isEntity(entity *Entity, out chan bool) {
	data := make(chan Entities)
	go n.retrieve(entity.Center, data)
	entities := <-data
	out <- entities.Contains(entity)
}

// split creates the children node for this node.
func (n *node) split() {
	n.children = append(n.children,
		n.new(NewBound(n.bound.Min.X, n.bound.Center.Y, n.bound.Center.X, n.bound.Max.Y)), // Top Left child node
		n.new(NewBound(n.bound.Center.X, n.bound.Center.Y, n.bound.Max.X, n.bound.Max.Y)), // Top Right child node
		n.new(NewBound(n.bound.Min.X, n.bound.Min.Y, n.bound.Center.X, n.bound.Center.Y)), // Bottom Left child node
		n.new(NewBound(n.bound.Center.X, n.bound.Min.Y, n.bound.Max.X, n.bound.Center.Y)), // Bottom Right child node
	)
}

func (n *node) MoveEntities(entities Entities, maxDepth uint16) {
	// loop through all entities to add them to there appropriate child node
	for i := range entities {
		// get the next node that the given entity fits in and insert it
		n.insert(entities[i], maxDepth)
	}

	// clear entities for branch node
	n.entities = n.entities[:0]
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
