package quadgo

import (
	"errors"
	"fmt"
)

// Entities is a list of entities.
type Entities []*Entity

// FindAndRemove finds and removes the given entity from the list of entities.
// remove returns the new list of entities and an error if the given entity can not be found in the list of entities.
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

// isIntersectPoint finds if a given point intersects any entities  in
// the list of entities. It returns a bool on an output chan for running on a
// secondary thread.
func (e Entities) isIntersectPoint(point Point, out chan<- bool) {
	// check if any entities returned intersect the given point
	for i := range e {
		// check for intersect
		if e[i].IsIntersectPoint(point) {
			out <- true
		}
	}
	out <- false
}

// isIntersectBound finds if a given bound intersects any entities  in
// the list of entities. It returns a bool on an output chan for running on a
// secondary thread.
func (e Entities) isIntersectBound(bound Bound, out chan<- bool) {
	// check if any entities returned intersect the given point
	for i := range e {
		// check for intersect
		if e[i].IsIntersectBound(bound) {
			out <- true
		}
	}
	out <- false
}

// isIntersectsPoint finds if a given point intersects any entities  in
// the list of entities. It returns a list of intersected entities
// on an output chan for running on a secondary thread.
func (e Entities) intersectsPoint(point Point, out chan<- Entities) {
	entities := make(Entities, 0, cap(e))

	// check if any entities returned intersect the given point and if they do add them to the return list
	for i := range e {
		// add to list if they intersect
		if e[i].IsIntersectPoint(point) {
			entities = append(entities, e[i])
		}
	}

	out <- entities
}

// isIntersectsBound finds if a given Bound intersects any entities  in
// the list of entities. It returns a list of intersected entities
// on an output chan for running on a secondary thread.
func (e Entities) intersectsBound(bound Bound, out chan<- Entities) {
	entities := make(Entities, 0, cap(e))

	// check if any entities returned intersect the given point and if they do add them to the return list
	for i := range e {
		// add to list if they intersect
		if e[i].IsIntersectBound(bound) {
			entities = append(entities, e[i])
		}
	}

	out <- entities
}

// Action is a function type that can be given to a entity to be executed later.
type Action func()

// Entity is the basic Entity stricture type for QuadGo.
//
// Entity holds the Bound information for an entity in the tree and also a list of interface{} which can hold
// any data that you would want to store in the entity.
type Entity struct {
	Bound
	Action
}

// NewEntity creates a new entity from the given min and max points and any given objects.
//
// The given objects can be any data that you want to hold with in the entity for the given bounds.
func NewEntity(minX, minY, maxX, maxY float64) *Entity {
	return &Entity{
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
		Bound:  NewBound(minX, minY, maxX, maxY),
		Action: action,
	}
}

// SetAction sets an entities action function.
func (e *Entity) SetAction(action Action) {
	e.Action = action
}

// IsEqual checks if the given entities bound is equal to this entities bound.
//
// IsEqual ignores the action function in comparison as you can not compare anonymous functions.
func (e *Entity) IsEqual(entity *Entity) bool {
	return e.Bound.IsEqual(entity.Bound)
}

func (e *Entity) String() string {
	return fmt.Sprintf("Bounds: %v Action: %v\n", e.Bound, e.Action)
}
