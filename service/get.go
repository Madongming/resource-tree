package service

import (
	"cache"
	"model"
)

func GetUserNode(userId int) (int64, error) {
	user, err := model.GetUser(userId)
	if err != nil {
		return 0, err
	} else if user == nil {
		return 0, nil
	}
	return int64(user.RootNode), nil
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

func GetGroupUsers(groupId int) ([]*User, error) {
	modelUsers, err := model.GetGroupUsers(groupId)
	if err != nil {
		return nil, err
	}
	return changeUsersModel2Service(modelUsers), nil
}

func GetUserPermission(userId, nodeId int) (int64, error) {
	permission, err := model.GetNodeUserPermission(userId, nodeId)
	if err != nil {
		return 0, err
	}

	return int64(permission.ReadWriteMask), nil
}

func GetGroupPermission(groupId, nodeId int) (int64, error) {
	permission, err := model.GetNodeGroupPermission(groupId, nodeId)
	if err != nil {
		return 0, err
	}

	return int64(permission.ReadWriteMask), nil
}

func GetTreeNodeUsers(nodeId int) ([]*User, error) {
	modelUsers, err := model.GetTreeNodeUsers(nodeId)
	if err != nil {
		return nil, err
	}
	users := changeUsersModel2Service(modelUsers)
	return users, nil
}

func GetTreeNodeGroups(nodeId int) ([]*Group, error) {
	modelGroups, err := model.GetTreeNodeGroups(nodeId)
	if err != nil {
		return nil, err
	}
	groups := changeGroupsModel2Service(modelGroups)
	return groups, nil
}
