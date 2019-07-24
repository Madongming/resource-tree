package cache

import (
	"model"
)

func (c *Cache) makeNodeArray2Tree(nodes []*model.ResourceTreeNode) error {
	// Get cahche the structrue from the model.
	c.changeModel2Resource(nodes)

	// 最大id + 1，做为索引长度
	indexLen := c.ReData[len(c.ReData)-1].ID + 1
	// 利用父ID建立树形索引
	// 第二个参数为了防止动态扩充索引切片
	err := c.makeArray2Tree(indexLen)
	if err != nil {
		return err
	}

	// 利用节点ID建立索引树节点的索引
	// 第二个参数为了防止动态扩充索引切片
	c.makeTreeIndex(indexLen)

	return nil
}

func (c *Cache) changeModel2Resource(nodes []*model.ResourceTreeNode) {
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for i := range nodes {
		c.ReData[i].ID = nodes[i].ID
		c.ReData[i].Parent = nodes[i].Parent
		c.ReData[i].Description = nodes[i].Description
		c.ReData[i].Level = nodes[i].Level
		c.ReData[i].Name = nodes[i].Name
		c.ReData[i].CnName = nodes[i].CnName
		c.ReData[i].Key = nodes[i].Key
	}
}

func (c *Cache) makeArray2Tree(indexLen int) error {
	if indexLen == 0 || c.ReData == nil || len(c.ReData) == 0 {
		return nil
	}

	parentIndex := make([]*ResourceTreeNode, indexLen)
	// 建立 父节点ID-> 子节点序列的映射，最多n个
	for i := 1; i < indexLen; i++ {
		parentIndex[c.ReData[i].Parent] = append(c.ReData[c.ReData[i].Parent],
			c.ReData[i])
	}

	// In database, must set root node's parent is 0.
	c.ReTree = new(ResourceTree)
	c.ReTree.Node = parentIndex[0]
	if c.ReTree.Node == nil {
		return ERR_ROOT_NODE_NOT_EXIST
	}

	makeTree(c.ReTree, parentIndex)
	return nil
}

func makeTree(root *ResourceTree, parentIndex [][]*ResourceTreeNode) {
	root.Childs = parentIndex[root.Node.ID]
	if root.Childs == nil || len(root.Childs) == 0 {
		return
	}

	for i := range root.Childs {
		makeTree(root.Childs[i], parentIndex)
	}
}

func (c *Cache) makeTreeIndex(indexLen int) {
	if c.ReTree == nil {
		return
	}
	c.ReIndex = make([]*ResourceTree, indexLen)
	makeIndex(c.ReTree, c.ReIndex)
	return index
}

func makeIndex(tree *ResourceTree, index []*ResourceTree) {
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
