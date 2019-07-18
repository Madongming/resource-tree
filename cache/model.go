package cache

var Cache *Resource

type Resource struct {
	Tree    *ResourceTree
	Index   []*ResourceTree
	Version int64
}

// For api response and in cache.
type ResourceTree struct {
	Node   *ResourceTreeNode `json:"node"`
	Childs []*ResourceTree   `json:"childs"`
}

type ResourceTreeNode struct {
	ID          int    `json:"id"`
	Parent      int    `json:"parent"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	EnName      string `json:"enName"`
	Key         string `json:"key"`
}
