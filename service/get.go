package service

import (
	"cache"
	"model"
)

func GetUserNode(userId int) (int, error) {
	user, err := model.GetUser(userId)
	if err != nil {
		return 0, err
	} else if user == nil {
		return 0, nil
	}
	return user.RootNode, nil
}

func GetTree(nodeId int) (*cache.ResourceTree, error) {
	if err := CasResource(); err != nil {
		return nil, err
	}

	tree, err := cache.FindTree(nodeId)
	if err != nil {
		return nil, nil
	}
	return tree, nil
}
