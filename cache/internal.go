package cache

import (
	"github.com/Madongming/resource-tree/model"
)

func newLRU(size int64) *LRU {
	lru := new(LRU)
	lru.Index = make(map[int]*CacheNode, size)

	// Make double link.
	dummy := new(CacheNode)
	p := dummy
	for i := 0; i < int(size); i++ {
		node := new(CacheNode)
		p.Next = node
		node.Pre = p
		p = p.Next
	}
	p.Next = dummy.Next
	dummy.Next.Pre = p
	lru.Data = &CacheList{
		UserCacheHead: dummy.Next,
		Size:          size,
	}

	return lru
}

// 设置缓存正在重新索引或完成索引
func setIsReData(i bool) {
	Mux.Lock()
	IsReData = i
	Mux.Unlock()
}

// 预先分配用于重新加载的数据对象
func (rn *ResourceNodeList) perMallocReData(dataLen int) {
	if dataLen == 0 {
		return
	}
	// 预分配树的对象，从对象池新取或者复用之前的对象

	if rn.ReData == nil {
		// 第一次更新
		rn.ReData = make([]*model.ResourceNode, dataLen)
		for i := range rn.ReData {
			// 从池子中取对象
			rn.ReData[i] = NodePool.Get().(*model.ResourceNode)
		}
	} else if len(rn.ReData) != dataLen {
		// 对象个数有变化
		// 将对象放回池子
		for i := range rn.ReData {
			NodePool.Put(rn.ReData[i])
		}
		rn.ReData = make([]*model.ResourceNode, dataLen)
		for i := range rn.Data {
			// 从池子中取对象
			rn.ReData[i] = NodePool.Get().(*model.ResourceNode)
		}
	}
	// 以上都未匹配，可以直接复用
}

func (tc *TreeCache) makeTree() error {
	// 最大id + 1，做为索引长度
	indexLen := ResourceNodes.ReData[len(ResourceNodes.ReData)-1].ID + 1
	// 利用父ID建立树形索引
	// 第二个参数为了防止动态扩充索引切片
	err := tc.makeArray2Tree(indexLen)
	if err != nil {
		return err
	}

	// 利用节点ID建立索引树节点的索引
	// 第二个参数为了防止动态扩充索引切片
	tc.makeTreeIndex(indexLen)

	return nil
}

func (rn *ResourceNodeList) changeModel2Resource(nodes []*model.DBResourceNode) {
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for i := range nodes {
		rn.ReData[i].ID = nodes[i].ID
		rn.ReData[i].Parent = nodes[i].Parent
		rn.ReData[i].Description = nodes[i].Description
		rn.ReData[i].Level = nodes[i].Level
		rn.ReData[i].Name = nodes[i].Name
		rn.ReData[i].CnName = nodes[i].CnName
		rn.ReData[i].Key = nodes[i].Key
	}
}

func (tc *TreeCache) makeArray2Tree(indexLen int) error {
	if indexLen == 0 ||
		ResourceNodes.ReData == nil ||
		len(ResourceNodes.ReData) == 0 {
		return nil
	}

	parentIndex := make([][]*model.ResourceNode, indexLen)
	// 建立 父节点ID-> 子节点序列的映射，最多n个
	for i := range ResourceNodes.ReData {
		parentIndex[ResourceNodes.
			ReData[i].
			Parent] = append(parentIndex[ResourceNodes.
			ReData[i].
			Parent],
			ResourceNodes.ReData[i])
	}

	// 虚拟Root
	tc.ReTree = new(model.Tree)
	tc.ReTree.Node = new(model.ResourceNode)
	tc.ReTree.Node.Parent = -1

	makeTree(tc.ReTree, parentIndex)
	return nil
}

func makeTree(root *model.Tree, parentIndex [][]*model.ResourceNode) {
	childs := parentIndex[root.Node.ID]
	for i := range childs {
		root.Childs = append(root.Childs,
			&model.Tree{
				Node: childs[i],
			})
	}
	if root.Childs == nil || len(root.Childs) == 0 {
		return
	}

	for i := range root.Childs {
		makeTree(root.Childs[i], parentIndex)
	}
}

func (tc *TreeCache) makeTreeIndex(indexLen int) {
	if tc.ReTree == nil {
		return
	}
	tc.ReIndex = make([]*model.Tree, indexLen)
	makeIndex(tc.ReTree, tc.ReIndex)
}

func makeIndex(tree *model.Tree, index []*model.Tree) {
	if tree == nil {
		return
	}

	if tree.Node == nil {
		return
	} else {
		index[tree.Node.ID] = tree
	}

	if tree.Childs == nil || len(tree.Childs) == 0 {
		return
	}

	for i := range tree.Childs {
		makeIndex(tree.Childs[i], index)
	}
}

// 交换重新索引的数据及旧数据
func (tc *TreeCache) swapReData() {
	Mux.Lock()
	tc.Index, tc.ReIndex = tc.ReIndex, tc.Index
	tc.Version, tc.ReVersion = tc.ReVersion, tc.Version
	Mux.Unlock()
}

func (rn *ResourceNodeList) swapReData() {
	Mux.Lock()
	rn.Data, rn.ReData = rn.ReData, rn.Data
	rn.Version, rn.ReVersion = rn.ReVersion, rn.Version
	Mux.Unlock()
}

func changeHeadPreAndUpNode(headPreNode,
	updatedNode *CacheNode) {
	if headPreNode == nil ||
		updatedNode == nil {
		return
	}
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

func newTreeByPermission(tree, newTree *model.Tree, permissionSet map[int]struct{}) {
	if _, found := permissionSet[tree.Node.ID]; found {
		// 找到有权限的节点，将其及子节点都加入树
		newTree.Childs = append(newTree.Childs,
			tree)
		return
	}
	if tree.Childs == nil || len(tree.Childs) == 0 {
		return
	}
	for i := range tree.Childs {
		newTreeByPermission(tree.Childs[i], newTree, permissionSet)
	}
}
