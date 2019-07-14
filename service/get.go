package service

import (
	"cache"
)

func GetTree(user *mode.User) (*cache.ResourceTree, error) {
	if err := CasResource(); err != nil {
		return nil, err
	}

	tree, err := cache.FindTree(user.RootNode)
	if err != nil {
		return nil, nil
	}
	return tree, nil
}
