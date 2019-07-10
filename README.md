# QuadGo [![GoDoc](https://godoc.org/github.com/Tskken/QuadGo?status.svg)](https://godoc.org/github.com/Tskken/QuadGo) ![Go Report Card](https://goreportcard.com/badge/github.com/Tskken/QuadGo) [![Build Status](https://travis-ci.org/Tskken/QuadGo.svg?branch=master)](https://travis-ci.org/Tskken/QuadGo)

QuadGO is a Quadtree implementation aimed at being used for video game collision detection.
The main goal of this library is to create an easy to use and easily extendable quadtree implementation
in Golang.

## Getting Started
To get QuadGo run  `go get github.com/Tskken/quadgo` in your command line of choice.
Then add it to any of your existing projects by adding `import "github.com/Tskken/quadgo"`.

## Tutorial

#### Creating the quadtree

First you need to create a QuadGo instance. For a basic tree using pre-definde presets you only need to use
the `tree := quadgo.New()` call. This will create and return a new quadtree with the pre-defined settings.

The defaults for this library are set to:
- Bounds: 1024x768
- Max entities per node: 10
- Max depth of the tree: 2

If you want to change the default settings you pass New() the options functions with there given arguments. For
example if you wanted to change the root bounds of the tree to be 1920x1080 you would do `quadgo.New(quadgo.SetBounds(1920, 1080))`

The current available options you can give to New() are:
- SetBounds(width, height float64)
- SetMaxEntities(maxEntities int)
- SetMaxDepth(maxDepth int)

In the future there may be more options added. If there are more you can just add them to the New() call with
no need to change your old code other wise.

The most common options most people will likely use is quadgo.SetBounds() as your game world will likely
not have a size of 1024x760, but this was a simple preset and allows for quick use if your bounds is in fact this size.

#### Added entities to the tree

The simplest way to add something to the tree is to call `tree.Insert(minX, minY, maxX, maxY, objects...)`. 
This will insert the given data in to the tree in what ever quadrant it fits in to. When you call Insert() it takes the min and max
bounds of the object you are trying to insert and any number of other "objects" that you may want to add to that entity. These
objects can be anything as it will be sorted as an array of interface{}. This can be used later as a way to do an action on something
when you retrieve an entity from the tree through something like `tree.IntersectsPoint()`.

If you already have some quadgo.Entity items that you want to insert in to the tree you can also just call `tree.InsertEntities(entities...)`.
This will insert any number of quadgo.Entity's in to the tree. If you do want to use this method then you will want to create your own entities.
This can be done by calling `quadgo.NewEntity(minX, minY, maxX, maxY, objects...)`. You could also just create the entity by hand as all structs are
exported for your convenience. Note though that the simplest way to create data is to just use the built in functions or to just insert data threw
Insert() as it creates the entity for you.

#### Removing entities from the tree

If you want to remove an entity from the tree you need to call `tree.Remove(entity)`. This will try and remove the given entity from the tree.
The given entity does not need to be the exact entity you are trying to remove but rather it just has to have the same data. This means if you want to remove
an entity from the node you just need to create an entity with the same data and call Remove() with that entity.

#### Retrieving entities from the tree

If you want to retrieve entities from the tree you just need to call ether `tree.RetrieveFromPoint(point)` or `tree.RetrieveFromBound(bound)`.
Both of these functions will return a list of entities that can be found at what ever leaf node the given point or bound can fit with in.

#### Checking for collisions

One of the key points of any quadtree is its collision detection, at least for games. You can check for collision with any entity with in the quadtree by
calling `tree.IsIntersectPoint(point)` or `tree.IsIntersectBound(bound)`. These two functions function the same way, just one checks for a point collision and one
checks if any point of the given bounds intersects anything with in the tree. These two functions return a boolean that will be true if any intersect
is found and false if not.

Another important function is the `tree.IntersectsPoint(point)` and `tree.IntersectsBound(bound)`. These two functions are almost the exact same as 
IsIntersectPoint and IsIntersectBound but rather then returning a boolean it returns a list of all entities that the given point or bound intersects with.

## License

[MIT](LICENSE)
