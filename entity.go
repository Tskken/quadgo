package QuadGo

// Entity is the basic Entity stricture type for QuadGo.
//
// Entity holds the bounds information for an entity in the tree and also a interface{} which can hold.
// any data that would be needed to be stored in the tree with the bounds.
// This could be used for things like actions on intersect with an entity or some extra information other then
// the basic bounds data.
type Entity struct {
	bounds Bounds
	object interface{}
}
