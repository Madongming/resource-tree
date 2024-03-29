package cache

import (
	"sync"

	"github.com/Madongming/resource-tree/model"
)

type ResourceNodeList struct {
	Data    []*model.ResourceNode
	Version int

	ReData    []*model.ResourceNode
	ReVersion int
}

type TreeCache struct {
	Tree    *model.Tree
	Index   []*model.Tree
	Size    int
	Version int

	ReTree    *model.Tree
	ReIndex   []*model.Tree
	ReSize    int
	ReVersion int
}

type LRU2 struct {
	Cache1  *LRU
	Cache2  *LRU
	mux     sync.Mutex
	Version int
}

type LRU struct {
	Index map[int]*CacheNode
	Data  *CacheList
	mux   sync.Mutex
}

// 环形双向链表
type CacheList struct {
	UserCacheHead *CacheNode
	Size          int64
}

// 链表节点
type CacheNode struct {
	Key  int
	Val  interface{}
	Pre  *CacheNode
	Next *CacheNode
}
