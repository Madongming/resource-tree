package cache

var Cache *Resource

type Resource struct {
	Tree    *ResourceTree
	Version int64
}

type ResourceTree struct {
	Node   *ResourceTreeNode
	Childs []*ResourceTree
}

type ResourceTreeNode struct {
	ID          int
	Parent      int
	Description string
	Level       int
	Name        string
	EnName      string
	Key         string
}
