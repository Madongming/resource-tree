package cache

import (
	"model"
)

func makeNodeArray2Tree(nodes []*model.ResourceTreeNode, version int64) (*Resource, error) {
	// 利用父ID建立树形索引
	tree, err := makeArray2Tree(nodes)
	if err != nil {
		return nil, err
	}

	// 利用节点ID建立索引
	maxNumOfNodeID := nodes[len(node)-1].ID
	index := make([]*ResourceTree, maxNumOfNodeID+1)
	makeTreeIndex(tree, index)

	return &Resource{
		Tree:    tree,
		Index:   index,
		Version: version,
	}, nil

}

func makeArray2Tree(nodes []*model.ResourceTreeNode) (*ResourceTree, error) {
	// 建立 父节点ID-> 子节点序列的映射，最多n个
	parentMap := make(map[int][]*ResourceTreeNode, len(nodes))
	for i := range nodes {
		parentMap[nodes[i].Parent] = append(
			parentMap[nodes[i].Parent], nodes[i])
	}

	// In database, must set root node's parent is -1.
	root, found := parentMap[-1]
	if !found {
		return nil, ERR_ROOT_NODE_NOT_EXIST
	}
	makeTree(root, parentMap)
	return root, nil
}

func makeTree(root *ResourceTreeNode, parentMap map[int][]*ResourceTreeNode) {
	var found bool
	root.Childs, found = parentMap[root.ID]
	if !found {
		return
	}

	for i := range root.Childs {
		makeTree(root.Childs[i], parentMap)
	}
}

func makeTreeIndex(tree *ResourceTree, index []*ResourceTree) {
	if tree == nil {
		return
	} else if tree.Node != nil {
		index[tree.Node.ID] = tree
	}
	if tree.Childs == nil || len(tree.Childs) == 0 {
		return
	}

	for i := range tree.Childs {
		makeTreeIndex(tree.Childs[i], index)
	}
}
