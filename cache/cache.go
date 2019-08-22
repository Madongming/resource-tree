package cache

import (
	"sync"

	. "github.com/Madongming/resource-tree/global"
	"github.com/Madongming/resource-tree/model"
	"github.com/Madongming/resource-tree/tools"
)

var (
	ResourceNodes *ResourceNodeList // 所有节点
	IsReData      bool              // 正在重新获取数据库中的数据并重建某个索引
	Mux           sync.Mutex

	Tree        *TreeCache // 缓存所有节点的树
	UserTreeLRU *LRU       // 缓存用户有权限节点的子树
	ResourceLRU *LRU2      // 缓存资源与其他资源的关联图

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
		UserTreeLRU = newLRU(Configs.UserCacheSize)
	}
	if ResourceLRU == nil {
		ResourceLRU = new(LRU2)
		ResourceLRU.Cache1 = newLRU(Configs.ResourceCacheSize)
		ResourceLRU.Cache2 = newLRU(Configs.GraphCacheSize)
	}
}

func (tc *TreeCache) Set(version int, data []*model.DBResourceNode) error {
	if IsReData {
		return ERR_RE_INDEX
	}

	// 设置为正在重新索引
	setIsReData(true)

	ResourceNodes.ReVersion = version
	tc.ReVersion = version
	tc.ReSize = len(data)
	// 预分配加载数据的对象
	ResourceNodes.perMallocReData(len(data))

	// 将从数据库取出的节点转成缓存里的节点.
	ResourceNodes.changeModel2Resource(data)
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

// 根据节点ID获取以该节点为root的子树
func (tc *TreeCache) GetTreeNode(nodeId int) (*model.Tree, error) {
	if nodeId > len(tc.Index) || tc.Index[nodeId] == nil {
		return nil, ERR_NODE_NOT_EXIST
	}
	return tc.Index[nodeId], nil
}

func (tc *TreeCache) Get() (*model.Tree, error) {
	return tc.Tree, nil
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

func (l2 *LRU2) Set(keys []int, value interface{}) {
	if keys == nil || len(keys) == 0 {
		return
	}

	l2.mux.Lock()
	middleKey := tools.GetRandInt()
	l2.Cache2.Set(middleKey, value)

	for i := range keys {
		l2.Cache1.Set(keys[i], middleKey)
	}
	l2.mux.Unlock()
}

func (l2 *LRU2) Get(key int) (interface{}, error) {
	v1, err := l2.Cache1.Get(key)
	if err != nil {
		return nil, err
	}
	key2, ok := v1.(int)
	if !ok {
		return nil, ERR_ASSERTION
	}
	return l2.Cache2.Get(key2)
}

func NewTreeByPermission(permissionSet map[int]struct{}) (*model.Tree, error) {
	newTree := new(model.Tree)
	newTree.Node = new(model.ResourceNode)
	newTree.Node.Parent = -1

	newTreeByPermission(Tree.Tree, newTree, permissionSet)
	return newTree, nil
}
