package cache

import (
	"model"
)

func makeNodeArray2Tree(nodes []*model.ResourceTreeNode, version int64) (*Resource, error) {
	tree, err := makeArray2Tree(nodes)
	if err != nil {
		return nil, err
	}
	index, err := makeTreeIndex(tree)
	return &Resource{
		Tree:    tree,
		Index:   index,
		Version: version}, nil

}

func makeArray2Tree(nodes []*model.ResourceTreeNode) (*ResourceTree, error) {
	// 当每个元素都有自己唯一的父节点，最多n个
	parentMap := make(map[int][]*ResourceTreeNode, len(nodes))
	for i := range nodes {
		parentMap[nodes[i].Parent] = append(
			parentMap[nodes[i].Parent], nodes[i])
	}
	root, found := parentMap[-1]
	if !found {
		return nil, ERR_ROOT_NODE_NOT_EXIST
	}
	makeTree(root, nodeMap, parentMap)
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

func makeTreeIndex(tree *ResourceTree) (map[int]*ResourceTree, error) {

}
