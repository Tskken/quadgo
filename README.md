# QuadGo [![GoDoc](https://godoc.org/github.com/Tskken/QuadGo?status.svg)](https://godoc.org/github.com/Tskken/QuadGo) [![Go Report Card](https://goreportcard.com/badge/github.com/Tskken/QuadGo)](https://goreportcard.com/report/github.com/Tskken/QuadGo)

QuadGO is a Quadtree implementation aimed at being used for video game collision detection.
The main goal of this library is to create an easy to use and easily extendable quadtree implementation
in Golang.

## Getting Started
To get QuadGo run  `go get github.com/Tskken/QuadGo` in your command line of choice.
Then add it to any of your existing projects by adding `import "github.com/Tskken/QuadGo".

## Tutorial

First you need to create a QuadGo instance.

```go
Quad := NewQuadGo(10, NewBounds(0, 0, 1024, 768))
```

The first value in NewQuadGo() sets the max number of entities a given node can have
before it splits. The Second value is the root bounds of the tree. This tipicly will be what 
ever your window and image size is.

#### Added data to the tree

To add data to the tree you just need to do 
```go
Quad.Insert(entity)
```

entity is anything that implements the Entity interface. If you do not want to have to write your
own Entity type you can just use QuadGo's Bounds data type with NewBounds().

#### Retrieving entities from the tree

If you want to find entities with in the tree you just need to do
```go
Quad.Retrieve(entity)
```

Quad.Retrieve() returns a list of all entities in the tree that can be found with in any
bounds that the given entity fits within.

#### Checking for collisions

If you want to check if an entity collides with anything with in the tree all you need to do is
```go
Quad.IsIntersect(entity)
```
This returns a boolean of true if entity collides with any other entity in the tree or false if
it does not.

#### Getting all collissions

Lastly if you want to get all the entitys that your entity collides with you just need to do
```go
Quad.Intersects(entity)
```
This returns a list of all entities that the given entity collided with.

## TODO's

- Threading
- Add removal of entities
- ...

## License

[MIT](LICENSE)
