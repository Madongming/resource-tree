package cache

import (
	"model"
)

func makeNodeArray2Tree(nodes []*model.ResourceTreeNode, version int64) (*Resource, error) {
	// Get cahche the structrue from the model.
	resourceTreeNodes := changeModel2Resource(nodes)

	// 最大id + 1，做为索引长度
	indexLen := resourceTreeNodes[len(resourceTreeNodes)-1].ID + 1
	// 利用父ID建立树形索引
	// 第二个参数为了防止动态扩充索引切片
	tree, err := makeArray2Tree(resourceTreeNodes, indexLen)
	if err != nil {
		return nil, err
	}

	// 利用节点ID建立索引树节点的索引
	// 第二个参数为了防止动态扩充索引切片
	index := makeTreeIndex(tree, indexLen)

	return &Resource{
		Tree:    tree,
		Index:   index,
		Version: version,
	}, nil

}

func changeModel2Resource(nodes []*model.ResourceTreeNode) []*ResourceTreeNode {
	if nodes == nil || len(nodes) == 0 {
		return []*ResourceTreeNode{}
	}
	results := make([]*ResourceTreeNode, len(nodes))
	for i := range nodes {
		results[i] = &ResourceTreeNode{
			ID:          nodes[i].ID,
			Parent:      nodes[i].Parent,
			Description: nodes[i].Description,
			Level:       nodes[i].Level,
			Name:        nodes[i].Name,
			CnName:      nodes[i].CnName,
			Key:         nodes[i].Key,
		}
	}
	results[0] = (*ResourceTreeNode)(nil)
	return results
}

func makeArray2Tree(resourceTreeNodes []*ResourceTreeNode, indexLen int) (*ResourceTree, error) {
	if resourceTreeNodes == nil || len(resourceTreeNodes) == 0 {
		return nil, nil
	}

	parentIndex := make([]*ResourceTreeNode, indexLen)
	// 建立 父节点ID-> 子节点序列的映射，最多n个
	for i := 1; i < indexLen; i++ {
		parentIndex[resourceTreeNodes[i].Parent] = append(parentIndex[resourceTreeNodes[i].Parent],
			resourceTreeNodes[i])
	}

	// In database, must set root node's parent is 0.
	root := parentIndex[0]
	if root == nil {
		return nil, ERR_ROOT_NODE_NOT_EXIST
	}

	makeTree(root, parentIndex)
	return root, nil
}

func makeTree(root *ResourceTreeNode, parentIndex [][]*ResourceTreeNode) {
	root.Childs = parentIndex[root.ID]
	if root.Childs == nil {
		return
	}

	for i := range root.Childs {
		makeTree(root.Childs[i], parentIndex)
	}
}

func makeTreeIndex(tree *ResourceTree, indexLen int) []*ResourceTree {
	if tree == nil {
		return
	}
	index := make([]*ResourceTree, index)
	makeIndex(tree, index)
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
