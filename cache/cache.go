package cache

import (
	. "global"
)

var (
	ResourceNodes *ResourceNodeList // 所有节点
	IsReData      bool              // 正在重新获取数据库中的数据并重建某个索引
	Mux           sync.Mutex

	Tree        *TreeCache // 缓存所有节点的树
	UserTreeLRU *LRU       // 缓存用户有权限节点的子树
	ResourceLRU *LRU       // 缓存一个节点的关联图

	NodePool = sync.Pool{
		New: func() interface{} {
			return &model.ResourceNode{}
		},
	}
)

// 初始化LRU
func init() {
	if ResourceNodes == nil {
		ResourceNodes = new(ResourceNodeList)
	}
	if Tree == nil {
		Tree = new(TreeCache)
	}
	if UserTreeLRU == nil {
		UserTreeLRU = newLRU(Configs.UserCacheHead)
	}
	if ResourceLRU == nil {
		ResourceLRU = newLRU(Config.ResourceCacheSize)
	}
}

func (tc *TreeCache) Set(version int, data []*model.DBResourceTreeNode) error {
	if tc.Resource.IsReData {
		return ERR_RE_INDEX
	}

	// 设置为正在重新索引
	setIsReData(true)

	ResourceNodes.ReVersion = version
	tc.ReVersion = version
	// 预分配加载数据的对象
	ResourceNodes.perMallocReData(len(data))

	// 将从数据库取出的节点转成缓存里的节点.
	ResourceNodes.changeModel2Resource(nodes)
	// 在重新索引的数据集上产生树和id索引
	if err := tc.makeTree(); err != nil {
		// 失败，提前结束
		setIsReData(false)
		return err
	}

	// 交换新据、旧数据集
	tc.swapReData()
	ResourceNodes.swapReData()

	// 重新索引结束
	setIsReData(false)
	return nil
}

func (tc *TreeCache) GetTreeNode(nodeId int) (*model.Tree, error) {
	if nodeId > len(tc.Index) || tc.Index[nodeId] == nil {
		return nil, ERR_NODE_NOT_EXIST
	}
	return tc.Index[nodeId], nil
}

func (tc *TreeCache) Get(uid ...int) (*model.Tree, error) {
	if len(uid) == 0 {
		return tc.Tree, nil
	}

	// 根据用户删除没有权限的节点
	return nil, nil
}

func (l *LRU) Set(key int, value interface{}) {
	// 保证线程安全
	l.mux.Lock()
	defer l.mux.Unlock()

	if v, found := l.Index[key]; found {
		// 缓存存在，更新
		v.Val = value
		// 调节顺序
		changeHeadPreAndUpNode(
			l.Data.UserCacheHead.Pre, v)
	} else {
		// 缓存不存在，更新的数据覆盖head前一个节点
		l.Data.UserCacheHead.Pre.Val = value
		if l.Data.UserCacheHead.Pre.Key != 0 {
			// 该节点已经被使用过，删除Index中的key
			delete(l.Index,
				l.Data.UserCacheHead.Pre.Key)
		}
		l.Index[key] = l.Data.UserCacheHead.Pre
		l.Data.UserCacheHead.Pre.Key = key
	}

	// 向前移动环形链表的头，让头节点始终保持
	// 最新的数据，前一个节点是最旧的数据
	l.Data.UserCacheHead = l.Data.UserCacheHead.Pre
}

func (l *LRU) Get(key int) (interface{}, error) {
	value, found := l.Index[key]
	if !found {
		return nil, ERR_CACHE_KEY_NOT_EXIST
	}
	// 调节顺序
	changeHeadPreAndUpNode(
		l.Data.UserCacheHead.Pre, value)

	// 向前移动环形链表的头，让头节点始终保持
	// 最新的数据，前一个节点是最旧的数据
	l.Data.UserCacheHead = l.Data.UserCacheHead.Pre
	return value.Val, nil
}

func NewTreeByPermission(permissionSet map[int]struct{}) (*model.Tree, error) {
	newTree := new(model.Tree)
	newTreeByPermission(Tree, newTree, permissionSet)
	return newTree, nil
}
