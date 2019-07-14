package model

import (
	. "global"
)

func GetResourceVersion() (int64, error) {
	return getCurrentVersion()
}

func GetAllNodes() ([]*ResourceTreeNode, error) {
	// Fetch all node of the tree.
	var resourceTreeNodes []*ResourceTreeNode
	if result := DB().
		Find(&resourceTreeNodes); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
	}
	return resourceTreeNodes, nil
}
