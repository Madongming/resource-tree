package cache

import (
	"sync"
)

var (
	Cache       *Resource
	Mux         sync.Mutex
	UserTreeLRU *LRU
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

type LRU struct {
	Index map[int]*CacheNode
	Data  *CacheList
	Mux   sync.Mutex
}

// 环形双向链表
type CacheList struct {
	UserCacheHead *CacheNode
	Size          int64
}

// 链表节点
type CacheNode struct {
	Val    *ResourceTree
	UserId int
	Pre    *CacheNode
	Next   *CacheNode
}
