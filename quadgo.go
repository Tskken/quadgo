// Copyright 2019 Tskken. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// Package quadgo provides a basic quad-tree implementation for
// game collision detection. It provides most standard read and write
// operations for a quad-tree. The main uses cases would be using quadgo
// with its provided quadgo.IsIntersect() functions to check for intersects
// with in your game space. You would insert objects (quadgo.Entity) types with ether
// quadgo.Insert(), quadgo.InsertWithAction(), or quadgo.InsertEntities().
//
// Note that all read operations with in this library are run concurrently but not safe with
// write operations. No mutex locks are provided with in this library so if you want to make safe
// writes and reads concurrently, you would have to set up your own mutex lock functions.
package quadgo

import (
	"errors"
)

// Option function type for setting the options of a new tree.
type Option func(*options)

// options struct which holds all the information for creating a new quad-tree with its given information.
type options struct {
	MaxEntities uint64
	MaxDepth    uint16
}

// defaultOptions for QuadGo
var defaultOption = options{
	MaxEntities: 10,
	MaxDepth:    5,
}

// SetMaxEntities sets the max number of entities per each node in the new tree.
func SetMaxEntities(maxEntities uint64) Option {
	return func(o *options) {
		o.MaxEntities = maxEntities
	}
}

// SetMaxDepth sets the max depth that the tree can split to.
func SetMaxDepth(maxDepth uint16) Option {
	return func(o *options) {
		o.MaxDepth = maxDepth
	}
}

// QuadGo - Base quad-tree data structure.
type QuadGo struct {
	*node

	maxDepth uint16
}

// New creates the basic QuadGo instance.
//
// New requires a width and a height but can also be given any number of other supported Option functions.
//
// Example:
//  basic - quadgo.New(800, 600)
//  with option - quadgo.New(800, 600, SetMaxDepth(5))
//
// QuadGo sets the New defaults for max depth to 5 and max entities to 10.
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

// Insert takes the desired min and max xy points for the inserted entity.
//
// Insert will insert the entity for the given bounds in to all leaf nodes that
// the given bounds intersects with. This can mean duplicate references if the given bound
// is large and can intersect many leaf nodes. These are Entity references which help save
// on memory use but be aware if you insert large objects it can hinder performance.
func (q *QuadGo) Insert(minX, minY, maxX, maxY float64) {
	q.insert(NewEntity(minX, minY, maxX, maxY), q.maxDepth)
}

// InsertWithAction takes the desired min and max xy points for the inserted entity and an Action function.
func (q *QuadGo) InsertWithAction(minX, minY, maxX, maxY float64, action Action) {
	q.insert(NewEntityWithAction(minX, minY, maxX, maxY, action), q.maxDepth)
}

// InsertEntities inserts any number of entities in the quad-tree.
//
// This will return an error if you do not give it any entities.
func (q *QuadGo) InsertEntities(entities ...*Entity) error {
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

// Remove removes the given Entity from the quad-tree.
//
// The given entity has to be the exact same as the one you want to delete. This function
// uses the Entities ID and a comparison with its Bounds to confirm that the found entity is
// in fact the entity to remove.
//
// This will return an error if the entity given was not found in the quad-tree.
func (q *QuadGo) Remove(entity *Entity) error {
	return q.remove(entity)
}

// Retrieve returns all entities from all nodes the given bounds intersects with.
// Retrieve excludes duplected entities.
//
// The return of this function is a <-channel of entities. This is due to
// the fact that all reads are run concurrently. If You want to just wait for this
// function to return and block the hole time you can just call it with
// `entities := <-quadgo.Retrieve(bound)`. This will block until values are returned
// on the Entities chan.
//
// If you want to run retrieve and then do actions before retrieving the data from
// the Entities chan you can just save the chan with
// `out := quadgo.Retrieve(bound)`. You can then later use Go's `entities := <- out`
// to block till the entities are returned from retrieve.
func (q *QuadGo) Retrieve(bound Bound) <-chan Entities {
	out := make(chan Entities)

	go func() {
		out <- q.retrieve(bound)
		close(out)
	}()

	return out
}

// IsEntity checks if a given entity exists within the tree.
// The given entity has to be an exact value of the entity you want
// to find. This means it has to have the same ID and Bounds.
// Action functions are not compared.
//
// The return of this function is a <-channel of bool. This is due to
// the fact that all reads are run concurrently. If You want to just wait for this
// function to return and block the hole time you can just call it with
// `is := <-quadgo.IsEntity(entity)`. This will block until a value is returned
// on the chan.
//
// If you want to run IsEntity() and then do actions before retrieving the data from
// the chan you can just save the chan with `out := quadgo.IsEntity(entity)`.
// you can then later use Go's `is := <-out` to block until the value is returned from isEntity.
func (q *QuadGo) IsEntity(entity *Entity) <-chan bool {
	out := make(chan bool)

	go func() {
		out <- q.isEntity(entity)
		close(out)
	}()

	return out
}

// IsIntersect take a bound and returns if that bound intersects any entity within the tree.
//
// The return of this function is a <-channel of bool. This is due to
// the fact that all reads are run concurrently. If You want to just wait for this
// function to return and block the hole time you can just call it with
// `is := <-quadgo.IsIntersect(bound)`. This will block until a value is returned
// on the chan.
//
// If you want to run isIntersect and then do actions before retrieving the data from
// the chan you can just save the chan with `out := quadgo.IsIntersect(bound)`.
// you can then later use Go's `is := <-out` to block until the value is returned from isIntersect.
func (q *QuadGo) IsIntersect(bound Bound) <-chan bool {
	out := make(chan bool)

	go func() {
		out <- q.retrieve(bound).isIntersect(bound)
		close(out)
	}()

	return out
}

// Intersects takes a bound and returns all entities that the given bound intersects with.
// If no entities were found it will return an empty list of Entities.
//
// The return of this function is a <-channel of Entities. This is due to
// the fact that all reads are run concurrently. If You want to just wait for this
// function to return and block the hole time you can just call it with
// `entities := <-quadgo.Intersects(bound)`. This will block until entities are returned
// on the chan.
//
// If you want to run intersects and then do actions before retrieving the data from
// the chan you can just save the chan with `out := quadgo.Intersects(bound)`.
// you can then later use Go's `entities := <-out` to block until the value is returned from intersects.
func (q *QuadGo) Intersects(bound Bound) <-chan Entities {
	out := make(chan Entities)

	go func() {
		out <- q.retrieve(bound).intersects(bound)
		close(out)
	}()

	return out
}

// list of nodes
type nodes []*node

// node is the container that holds the branch and leaf data for the tree.
type node struct {
	parent   *node
	bound    Bound
	entities Entities
	children nodes
	depth    uint16
}

// new creates a new node instance for a given bounds taking the member node as its parent.
func (n *node) new(bound Bound) *node {
	return &node{
		parent:   n,
		bound:    bound,
		entities: make(Entities, 0, cap(n.entities)),
		children: make(nodes, 0, 4),
		depth:    n.depth + 1,
	}
}

// retrieve finds all of the entities with in a quadrant that the given point fits with in.
func (n *node) retrieve(bound Bound) (entities Entities) {
	// check if you are at a leaf node
	if len(n.children) > 0 {
		// get all child nodes the given bounds intersects
		nodes := n.getQuadrant(bound)
		// panic if no nodes were returned. This is a fatel error but should never happen in any normal situation.
		if len(nodes) == 0 {
			panic(errors.New("could not find a node to retrive from in node.retrieve().getQuadrent()"))
		}

		// recursive call to retrieve all entities from children nodes found from getQuadrent().
		for i := range nodes {
			ents := nodes[i].retrieve(bound)
			for i := range ents {
				if !entities.Contains(ents[i]) {
					entities = append(entities, ents[i])
				}
			}
		}
		return
	}

	// return entities from leaf
	return n.entities
}

// insert inserts a given entity in to the quad-tree.
func (n *node) insert(entity *Entity, maxDepth uint16) {
	// check if you are on a leaf node
	if len(n.children) > 0 {
		// get all child nodes the given bounds intersects
		nodes := n.getQuadrant(entity.Bound)
		// panic if no nodes were returned. This is a fatel error but should never happen in any normal situation.
		if len(nodes) == 0 {
			panic(errors.New("could not find a node to insert in to from node.insert().getQuadrent()"))
		}

		// recersive insert for all nodes found.
		for i := range nodes {
			nodes[i].insert(entity, maxDepth)
		}
		return
	}

	// check if a split is needed
	if len(n.entities)+1 > cap(n.entities) && n.depth < maxDepth {
		// split node in to child nodes
		n.split()

		// move this nodes entities to the children nodes
		n.moveEntities(append(n.entities, entity), maxDepth)
		return
	}

	// add Entity to node
	n.entities = append(n.entities, entity)
}

// remove removes the given Entity from the quadtree.
func (n *node) remove(entity *Entity) error {
	// check if we are on a leaf node
	if len(n.children) > 0 {
		// get all child nodes the given bounds intersects
		nodes := n.getQuadrant(entity.Bound)
		// panic if no nodes were returned. This is a fatel error but should never happen in any normal situation.
		if len(nodes) == 0 {
			panic(errors.New("could not find a node to remove in to from node.remove().getQuadrent()"))
		}

		// recersive call for all nodes found to remove from
		for i := range nodes {
			err := nodes[i].remove(entity)
			if err != nil {
				return err
			}
		}

		// collapse if needed if you are not at root
		if n.parent != nil {
			n.parent.collapse()
		}

		return nil
	}

	// find given entity in list of entities in node and return a list with that entity removed
	// returns an error if entity was not found
	entities, err := n.entities.FindAndRemove(entity)
	if err != nil {
		return err
	}

	// replace old version fo entities with new list of entities with given entity removed
	n.entities = entities

	// collapse if needed if you are not at root
	if n.parent != nil {
		n.parent.collapse()
	}

	return nil
}

// collapse takes all entities from the children nodes and moves them to the parent and then removes the children.
func (n *node) collapse() {
	// create an Entity array to coppy the entities to
	entities := make(Entities, 0, cap(n.entities))

	// cycle through children to find all non duplecet entities
	for i := range n.children {
		for _, ent := range n.children[i].entities {
			if !entities.Contains(ent) {
				entities = append(entities, ent)
			}
		}
	}

	// check if collapse is needed
	if len(entities) <= cap(n.entities) {
		// set parent entities to list of non duplecet entities
		n.entities = entities

		// clear children
		n.children = n.children[:0]
	}
}

// isEntity returns if a given entity exists in the tree.
func (n *node) isEntity(entity *Entity) bool {
	// check if you are at a leaf
	if len(n.children) > 0 {
		// get all child nodes the given bounds intersects
		nodes := n.getQuadrant(entity.Bound)
		// panic if no nodes were returned. This is a fatel error but should never happen in any normal situation.
		if len(nodes) == 0 {
			panic(errors.New("could not find a node for node.isEntity().getQuadrent()"))
		}

		// recersive call to check found nodes
		for i := range nodes {
			if nodes[i].isEntity(entity) {
				return true
			}
		}

		return false
	}

	// return if given entity is contained in entities list
	return n.entities.Contains(entity)
}

// split creates the children node for this node.
func (n *node) split() {
	n.children = append(n.children,
		n.new(NewBound(n.bound.Min.X, n.bound.Min.Y, n.bound.Center.X, n.bound.Center.Y)), // Top Left child node
		n.new(NewBound(n.bound.Center.X, n.bound.Min.Y, n.bound.Max.X, n.bound.Center.Y)), // Top Right child node
		n.new(NewBound(n.bound.Min.X, n.bound.Center.Y, n.bound.Center.X, n.bound.Max.Y)), // Bottom Left child node
		n.new(NewBound(n.bound.Center.X, n.bound.Center.Y, n.bound.Max.X, n.bound.Max.Y)), // Bottom Right child node
	)
}

// moveEntities moves the given entities to the children nodes of this node
func (n *node) moveEntities(entities Entities, maxDepth uint16) {
	// loop through all entities to add them to there appropriate child node
	for _, e := range entities {
		// get the next node that the given entity fits in and insert it
		n.insert(e, maxDepth)
	}

	// clear entities for branch node
	n.entities = n.entities[:0]
}

// getQuadrant returns the children nodes the given bound intersects with
func (n *node) getQuadrant(bound Bound) (nodes nodes) {
	for i := range n.children {
		if n.children[i].bound.IsIntersect(bound) {
			nodes = append(nodes, n.children[i])
		}
	}
	return
}
