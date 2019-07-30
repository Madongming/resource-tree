package cache

import (
	"sync"

	. "global"
	"model"
)

func init() {
	if Cache == nil {
		fetchCache()
	}
	if LRU == nil {
		initLRU()
	}
}

func CasResource() (bool, error) {
	version, err := model.GetResourceVersion()
	if err != nil {
		return false, err
	} else if Cache != nil && (Cache.IsReData || Cache.Version == version) {
		return false, nil
	}
	if err := fetchCache(version); err != nil {
		return false, err
	}
	return true, nil
}

func fetchCache(versions ...int64) error {
	if len(versions) == 0 {
		// 缓存中没有数据
		var err error
		Cache.ReVersion, err = model.GetResourceVersion()
		if err != nil {
			return err
		}
	} else {
		Cache.ReVersion = versions[0]
	}

	// 从数据库中获取全部节点数据
	allNodes, err := model.GetAllNodes()
	if err != nil {
		return err
	} else if len(allNodes) == 0 {
		return nil
	}

	// 设置为正在重新索引
	Cache.setIsReData(true)

	// 预分配加载数据的对象
	Cache.perMallocReData(len(allNodes))

	// 在重新索引的数据集上产生树和id索引
	if err := Cache.makeNodeArray2Tree(allNodes); err != nil {
		// 失败，提前结束
		Cache.setIsReData(false)
		return err
	}

	// 交换新据、旧数据集
	Cache.swapReData()

	// 重新索引结束
	Cache.setIsReData(false)
	return nil
}

// 设置缓存正在重新索引或完成索引
func (c *Cache) setIsReData(i bool) {
	Mux.Lock()
	c.IsReData = i
	Mux.Unlock()
}

// 预先分配用于重新加载的数据对象
func (c *Cache) perMallocReData(dataLen int) {
	if dataLen == 0 {
		return
	}
	// 预分配树的对象，从对象池新取或者复用之前的对象
	if Cache.ReData == nil {
		// 第一次更新
		Cache.ReData = make([]*ResourceTreeNode, dataLen)
		for i := range Cache.ReData {
			// 从池子中取对象
			Cache.ReData[i] = TreeNodePool.Get().(*ResourceTreeNode)
		}
	} else if len(Cache.ReData) != dataLen {
		// 对象个数有变化
		Cache.IsReData = true
		// 将对象放回池子
		for i := range Cache.ReData {
			TreeNodePool.Put(Cache.ReData[i])
		}
		Cache.ReData = make([]*ResourceTreeNode, dataLen)
		for i := range Cache.Data {
			// 从池子中取对象
			Cache.Data[i] = TreeNodePool.Get().(*ResourceTreeNode)
		}
	}
	// 以上都未匹配，可以直接复用
}

// 交换重新索引的数据及旧数据
func (c *Cache) swapReData() {
	Mux.Lock()
	c.Data, c.ReData = c.ReData, c.Data
	c.Index, c.ReIndex = c.ReIndex, c.Index
	c.Version, c.ReVersion = c.ReVersion, c.Version
	Mux.Unlock()
}

// 初始化LRU
func initLRU() {
	UserTreeLRU = new(LRU)
	UserTreeLRU.Index = make(map[int]*CacheNode, Configs.UserCacheSize)

	// Make double link.
	dummy := new(CacheNode)
	p := dummy
	for i := 0; i < Configs.UserCacheSize; i++ {
		node := new(CacheNode)
		p.Next = node
		node.Pre = p
		p = p.Next
	}
	p.Next = dummy.Next
	dummy.Next.Pre = p
	UserTreeLRU.Data = &CacheList{
		UserCacheHead: dummy.Next,
		Size:          Configs.UserCacheSize,
		Mux:           sync.Mutex{},
	}
}

func (ul *UserTreeLRU) Set(userId int, tree *ResourceTree) {
	// 保证线程安全
	ul.Mux.Lock()

	if v, found := ul.Index[userId]; found {
		// 缓存存在，更新
		v.Val = tree
		// 调节顺序
		changeHeadPreAndUpNode(
			ul.Data.UserCacheHead.Pre, v)
	} else {
		// 缓存不存在，更新的数据覆盖head前一个节点
		ul.Data.UserCacheHead.Pre.Val = tree
		if ul.Data.UserCacheHead.Pre.UserId != 0 {
			// 该节点已经被使用过，删除Index中的key
			delete(ul.Index,
				ul.Data.UserCacheHead.Pre.UserId)
		}
		ul.Index[userId] = ul.Data.UserCacheHead.Pre
		ul.Data.UserCacheHead.Pre.UserId = userId
	}

	// 向前移动环形链表的头，让头节点始终保持
	// 最新的数据，前一个节点是最旧的数据
	ul.Data.UserCacheHead = ul.Data.UserCacheHead.Pre

	ul.Mux.Unlock()
}

func (ul *UserTreeLRU) changeHeadPreAndUpNode(headPreNode,
	updatedNode *CacheNode) {
	headPreNodePre := headPreNode.Pre
	headPreNode.Pre = updatedNode.Pre
	updatedNode.Pre.Next = headPreNode

	headPreNodeNext := headPreNode.Next
	headPreNode.Next = updatedNode.Next
	updatedNode.Next.Pre = headPreNode

	updatedNode.Next = headPreNodeNext
	updatedNode.Pre = headPreNodePre
	headPreNodePre.Next = updatedNode
	headPreNodeNext.Pre = updatedNode
}
