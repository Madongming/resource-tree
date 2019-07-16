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
	if err := updateUsersOfGroupRootNode(groupId); err != nil {
		return err
	}
	return nil
}

func UpdateUserToNode(userId interface{}, nodeId int, permissions ...int) error {
	// ...

	if err := updateUserRootNode(userId); err != nil {
		return err
	}
	return nil
}

func UpdateGroupToNode(groupId interface{}, nodeId int, permissions ...int) error {
	//...

	if err := updateUsersOfGroupRootNode(groupId); err != nil {
		return err
	}
	return nil
}

func UpdateUserPermissionsToNode(userId, nodeId interface{}, permissions int) error {

}

func UpdateGroupPermissionsToNode(groupId, nodeId interface{}, permissions int) error {

}

func updateUserRootNode(userId interface{}) error {

}

func updateUsersOfGroupRootNode(groupId interface{}) error {

}
