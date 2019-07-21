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
	return model.AddUserToGroup(userId, groupId)
}

func AddUserToNode(userId interface{}, nodeId int, permissions ...int) error {
	// Default permission is 1.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	return model.AddUserToNode(userId, nodeId, permission)
}

func AddGroupToNode(groupId interface{}, nodeId int, permissions ...int) error {
	// Default permission is 1.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	return model.AddGroupToNode(groupId, nodeId, permission)
}
