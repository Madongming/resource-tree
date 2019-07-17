package service

import (
	"cache"
	. "global"
	"model"
)

// Add user when user first login.
func AddUser(name, cnName string) error {
	user, err := GetUserByName(name)
	if user != nil {
		return nil
	} else if err != nil {
		return err
	}
	return model.CreateUser(name, cnName)
}

// Create a group.
func CreateGroup(name, CnName string) error {
	return model.CreateGroup(name, cnName)
}

func AddUserToGroup(userId, groupId interface{}) error {
	if err := model.AddUserToGroup(userId, groupId); err != nil {
		return err
	}

	// Because, user's group is changed, it's perssion is maybe changed too.
	if err := updateUsersOfGroupRootNode(groupId); err != nil {
		return err
	}
	return nil
}

func AddUserToNode(userId interface{}, nodeId int, permissions ...int) error {
	// Default permission is readonly.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	if err := model.AddUserToNode(userId, nodeId, permission); err != nil {
		return err
	}

	// Update the root node of the user.
	if err := updateUserRootNode(userId); err != nil {
		return err
	}
	return nil
}

func AddGroupToNode(groupId interface{}, nodeId int, permissions ...int) error {
	// Default permission is readonly.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	if err := model.AddGroupToNode(groupId, nodeId, permission); err != nil {
		return err
	}

	// The permission of group is changed, the user's root node of the group is maybe changed too.
	if err := updateUsersOfGroupRootNode(groupId); err != nil {
		return err
	}
	return nil
}
