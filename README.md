# QuadGo [![GoDoc](https://godoc.org/github.com/Tskken/QuadGo?status.svg)](https://godoc.org/github.com/Tskken/QuadGo) ![Go Report Card](https://goreportcard.com/badge/github.com/Tskken/QuadGo) [![Build Status](https://travis-ci.org/Tskken/QuadGo.svg?branch=master)](https://travis-ci.org/Tskken/QuadGo)
 
QuadGO is a Quad-tree implementation aimed at being used for video game collision detection.
The main goal of this library is to create an easy to use and easily extendable quadtree implementation in Golang. It also attempts to tackle concerincy for all read operations, making use of Go's channel structures.
 
# Getting Started
To get QuadGo run  `go get github.com/Tskken/quadgo` in your command line of choice.
Then add it to any of your existing projects by adding `import "github.com/Tskken/quadgo"`.
 
# Tutorial
 
## Creating the quad-tree
 
To create your instance of QuadGo you have to call the New() function. This will create a new instance of QuadGo and return its reference to be used.
 
For example:
 
```go
    // create a basic instance with a given width and height
    tree := quadgo.New(width, height)
 
    // use quadgo Option to change quadgo defaults for a trea
    tree := quadgo.New(width, height, SetMaxDepth(depth))
```
 
QuadGo uses an Option's system for creation to make the new call both easy to use and easy to expand on if new options need to be added in the future. An Option is just a function type which changes the setting of the tree.
 
The current support options for quadgo.New() are:
- SetMaxEntities(uint64)
- SetMaxDepth(uint16)
 
The values are set to uint to enforce non negative value for SetMaxEntities and SetMaxDepth as you can not have a negative number of entities or depth of a tree.
 
The current defaults for the quad-tree are set to:
- Max entities per node: 10
- Max depth of the tree: 5
 
Note that New() can take any number of Option’s so you can pass more than one into the function at a time.
 
For example:
```go
    // create a tree with your own max entities and depth settings
    tree := quadgo.New(
        width, 
        height, 
        SetMaxEntities(maxEntities), 
        SetMaxDepth(maxDepth)
    )
```
 
## Added entities to the tree
 
By far the simplest way to insert any data into the tree is threw the quadgo.Insert() function. This function takes the min and max x and y positions for a new entity and creates and inserts it into the tree for you.
 
Example:
```go
    // insert an entity with a bounds of min:0, 0, max: 50, 50
    tree.Insert(0, 0, 50, 50)
```
 
This function as stated creates a new entity with the given bounds and inserts it into the tree. If you note through reading the godocs an entity also has an Action member which can be set with other function as shown in the next example. tree.Insert() sets Action to nil by default.
 
If you want to set the Action member for the inserted entity the easiest way would be to instead use InsertWithAction(). This function takes the bounds of the entity and a function which will be set as the Action function when creating the entity.
 
Example:
```go
    // insert an entity with the given bounds and an action function
    tree.InsertWithAction(0, 0, 50, 50, func(){
        fmt.Println("this is an action function print from entity")
    })
```
 
Additionally you can create your own entity or entities and insert them directly into the tree. To do this you use InsertEntities() which takes any number of Entity's as a variadic argument.
 
Example:
```go
    // create a list of entities 
    ...
 
    // insert some list of entities created prior to the tree
    err := tree.InsertEntities(entities)
    if err != nil {
        panic(err)
    }
```
 
Note that InsertEntities() does return an error. Because this function takes a variadic argument it will return an error if you call it with no entities. If you are sure there will be entities given to InsertEntities() then you can just ignore the error as no other part of the function will error.
 
## Removing entities from the tree
 
To remove entities from the tree you need to use quadgo.Remove(). This function will remove the given entity from the tree and if needed collapse any leafs to save memory space and clean up the tree.
 
Example:
```go
    // remove the entity from the tree
    tree.Remove(entity)
```
 
One important thing to understand is that the entity given to remove has to be an exact reference to the entity you want to remove. The check to find the entity uses the given entities ID and Bound. So this means you have to either already have or find the entity you want to remove from the tree first before you can remove it. This does make Remove a little less optimized in practice as you need to either keep copies of entity references or you need to find the entity first and then find it again in the tree and remove it.
 
#### Retrieving entities from the tree
 
To find entities in the tree you need to use quadgo.Retrieve(). This function takes a bounds to use to search the tree and will return all entities from nodes that that given entity intersects with.
 
Example:
```go
    // get all entities from nodes that bound intersects
    entities := <-tree.Retrieve(bound)
```
 
As a note if you see the `<-tree.` part of the code, this is because Retrieve()
returns a channel of Entities not just Entities. This is because all read operations in QuadGo run concurrently. In this case above we are calling tree.Retreive() and blocking till some amount of entities are returned on the output channel of Entities. We then assign that value to entities to use. Because of the fact that this function is concurrent that means we can also run the Retrieve() function and not try and receive from the channel till later similar to how a time.Timer() works in the Go standard library.
 
Example:
```go
    // get all entities for the given bounds but not receive
    // from the channel till later.
    out := tree.Retrieve(bound)
 
    // do some other stuff that doesn't need the value from Retrieve()
    ...
 
    // get the entities from the channel or block if needed
    entities := <-out
```
 
As you will see in the next section, and as stated earlier, all read operations for the tree are run concurrently. This means all read operations return a similar structure of a channel type to be used however is best for that instance.
 
One thing to note though is given that these operations are run concurrently, QuadGo does not have any native syncing functionality. This means that even though read operations are run concurrently, all write and update actions are not and could change data well a read operation is executing. This means you have to make sure that you are not using Insert() or Remove() well one of the read operations are running or you may get some unexpected behavior.
 
 
## Checking for collisions
 
One of the most important parts of a quad-tree made for collision detection in video games is the functions used to check for collisions. This takes the form of two functions with in QuadGo, the IsIntersect() and Intersects() functions. As the names kind of imply the IsIntersect() function checks if a given bounds intersects any entity with in the tree. The Intersects() function returns all entities tht the given bounds intersects with in the tree.
 
Example:
```go
    // check for collision with a bounds
    if <-tree.IsIntersect(bound) {
        // do something on intersect case
        ...
    }
 
    // get all entities the bounds intersects with
    entities := <-tree.Intersects(bound)
```
 
As you can see, and as said in the prior section, the collision check functions are also considered read functions and in turn run concurrently. This means they also return channels and can be used the same way Retrieve() was used up above.
 
Additional these functions run in to the same possible issues as Retrieve() as they are not safe with other write functions and could create unexpected behavior.
 
## Other useful functions
 
There is one other possibly useful function provided by QuadGo. This is the IsEntity() function. This function checks to see if the given entity exists with in the tree. Similery with Remove() the given entity has to be the same reference to the entity you are trying to find, but this could be useful if you want to check to make sure an entity was removed from the tree or to check to see if an entity exists with in the tree and if not add it back in.
 
Example:
```go
    // check if the entity exists with in the tree
    if !<-tree.IsEntity(entity) {
        // do something if entity didn't exist
        ...
    }
```
 
This function is also a read function so its run concurrently. So make sure if you are using it to check if an entity exists and if not add it back in to the tree, to make sure to block from the channel returned from IsEntity() before you do an insert or remove.
 
# Current TODO
- Re-write benchmarks for the library.
- More tests for possible edge cases.
 
## License
 
[MIT](LICENSE)
