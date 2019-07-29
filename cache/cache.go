package cache

import (
	"sync"

	"model"
)

func init() {
	if Cache == nil {
		fetchCache()
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
