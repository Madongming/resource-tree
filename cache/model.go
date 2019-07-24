package cache

import (
	"sync"
)

var (
	Cache *Resource
	Mux   sync.Mutex
)

var TreeNodePool = sync.Pool{
	New: func() interface{} {
		return &ResourceTreeNode{}
	},
}

type Resource struct {
	Tree    *ResourceTree
	Index   []*ResourceTree
	Version int64
	Data    []*ResourceTreeNode

	ReTree    *ResourceTree
	ReIndex   []*ResourceTree
	ReData    []*ResourceTreeNode
	ReVersion int64

	IsReData bool
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
	Name        string `json:"Name"`
	CnName      string `json:"cnName"`
	Key         string `json:"key"`
}
