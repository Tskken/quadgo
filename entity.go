// Copyright 2019 Tskken. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package quadgo

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Entities is a list of Entity's.
type Entities []*Entity

// FindAndRemove finds and removes the given entity from the list of entities.
// returns the new list of entities and an error if the given entity can not be found in the list of entities.
func (e Entities) FindAndRemove(entity *Entity) (Entities, error) {
	// check the entities in leaf for given entity
	for i := range e {
		// check if given entity is the same as nodes entity
		if e[i].IsEqual(entity) {
			// check if removal would make the leaf have no entities
			if len(e) == 1 {
				// set node entities to an empty slice
				e = e[:0]
			} else if len(e) == i+1 {
				// remove last entity from node
				e = e[:i]
			} else {
				// remove entity from node
				e = append(e[:i], e[i+1:]...)
			}
			return e, nil
		}
	}

	return nil, errors.New("could not find entity in tree to remove")
}

// Contains checks if the given entity exists with in the list of entities.
func (e Entities) Contains(entity *Entity) bool {
	// check each entity for if it is equal to given entity
	for i := range e {
		// check if given Entity equals given entity
		if e[i].IsEqual(entity) {
			return true
		}
	}
	return false
}

// isIntersectBound finds if a given bound intersects any entities  in
// the list of entities. It returns a bool on an output chan for running on a
// secondary thread.
func (e Entities) isIntersect(bound Bound) bool {
	// check if any entities returned intersect the given point
	for i := range e {
		// check for intersect
		if e[i].IsIntersect(bound) {
			return true
		}
	}
	return false
}

// isIntersectsBound finds if a given Bound intersects any entities  in
// the list of entities. It returns a list of intersected entities
// on an output chan for running on a secondary thread.
func (e Entities) intersects(bound Bound) (entities Entities) {
	// check if any entities returned intersect the given point and if they do add them to the return list
	for i := range e {
		// add to list if they intersect
		if e[i].IsIntersect(bound) {
			entities = append(entities, e[i])
		}
	}
	return
}

// Action is a function type that can be given to an entity to be executed later.
type Action func()

// Entity is the Entity structure type for QuadGo.
//
// Entity holds the Bound information for an entity in the tree and an Action function as a closer
// style function type which can store a function to use later. Entity also holds an ID which is
// by default a random uint64 value that is used to be able to accurately compare
// entities with IsEntity()
type Entity struct {
	ID uint64
	Bound
	Action
}

// NewEntity creates a new entity from the given min and max points.
//
// The ID for any given entity created will be default set to a random uint64 value seeded at creation
// time with time.Now().UnixNano(). If you want to set an ID you self just change the ID after creation.
func NewEntity(minX, minY, maxX, maxY float64) *Entity {
	return &Entity{
		ID:     rand.New(rand.NewSource(time.Now().UnixNano())).Uint64(),
		Bound:  NewBound(minX, minY, maxX, maxY),
		Action: nil,
	}
}

// NewEntityWithAction creates a new entity with the given min and max x and y positions of its bounds
// along with an Action function.
//
// Example:
//	quadgo.NewEntityWithAction(0, 0, 50, 50, func(){
//		fmt.Println("hello from an action")
//	})
func NewEntityWithAction(minX, minY, maxX, maxY float64, action Action) *Entity {
	return &Entity{
		ID:     rand.New(rand.NewSource(time.Now().UnixNano())).Uint64(),
		Bound:  NewBound(minX, minY, maxX, maxY),
		Action: action,
	}
}

// SetAction sets an entities action function.
func (e *Entity) SetAction(action Action) {
	e.Action = action
}

// IsEqual checks if the ID and bound of the entity is the same.
func (e *Entity) IsEqual(entity *Entity) bool {
	return (e.ID == entity.ID) && e.Bound.IsEqual(entity.Bound)
}

func (e *Entity) String() string {
	return fmt.Sprintf("ID: %v, Bounds: %v Action: %v\n", e.ID, e.Bound, e.Action)
}
