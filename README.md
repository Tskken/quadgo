# QuadGo [![GoDoc](https://godoc.org/github.com/Tskken/QuadGo?status.svg)](https://godoc.org/github.com/Tskken/QuadGo) [![Go Report Card](https://goreportcard.com/badge/github.com/Tskken/QuadGo)](https://goreportcard.com/report/github.com/Tskken/QuadGo) [![Build Status](https://travis-ci.org/Tskken/QuadGo.svg?branch=master)](https://travis-ci.org/Tskken/QuadGo)

QuadGO is a Quadtree implementation aimed at being used for video game collision detection.
The main goal of this library is to create an easy to use and easily extendable quadtree implementation
in Golang.

## Getting Started
To get QuadGo run  `go get github.com/Tskken/QuadGo` in your command line of choice.
Then add it to any of your existing projects by adding `import "github.com/Tskken/QuadGo".

## Tutorial

First you need to create a QuadGo instance.

```go
tree, err := NewQuadGo(10, 1024, 768)
```

The first value in NewQuadGo() sets the max number of entities a given node can have
before it splits. This number has to be greater then 0 and will return an error otherwise.
The second two values are the width and height of the screen or game world. This will be the maximum x and y
of the root of the tree. That means every object that is put in to the tree must fit with in these bounds.

#### Added entities to the tree

To add an entity to the tree you just need to do is call Insert() like so

```go
tree.Insert(bounds, object)
```

Insert() takes a Bounds structure type and an object as an interface. Anything can be put in
as an object and it should be used as a way to store extra data in the tree. Just note that
because you can put anything in this, if you put some large amount of data in it will make
the tree take up large amounts of space and may slow down search and removal speeds along with
just taking up a lot of memory.

#### Removing entities from the tree

To remove an entity you need to call Remove().

```go
tree.Remove(Entity)
```

Remove() will try and remove the given entity from the tree. Entity structure type
holds the bounds for that given entity and the object that it holds.

#### Retrieving entities from the tree

If you want to find entities with in the tree you just need to call Retrieve()

```go
tree.Retrieve(bounds)
```

Retrieve() returns a list of all entities that are with in the smallest leaf that the given
bounds can fit in. This can be the easiest way to get a list of entities that reside with in
an area of the screen.

#### Checking for collisions

If you want to check if a bounds collides with anything with in the tree you need to use IsIntersect()

```go
tree.IsIntersect(bounds)
```

This returns a true or false value for if the given bounds has intersected anything within
the tree. This function will be the fastest way to find a collision if your using this for
collision detection as it returns the second it finds a collision.

#### Getting all collision's

Lastly if you want to get a list of all collision's you need to call Intersects()

```go
tree.Intersects(bounds)
```

This functions almost identical to IsIntersects but instead of returning a true or false
it returns a list of all Entities that have been found to collide with the given bounds.
This will be best used if you have stored some action or data in the object field of 

## TODO's

- Threading (Tentitive may not be posible or worth doing. Currently under heavy reserch and testing)
- ...

## License

[MIT](LICENSE)
