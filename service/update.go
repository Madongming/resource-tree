package service

import (
	"model"
)

func UpdateUserPermissionsToNode(userId, nodeId interface{}, permission int) error {
	if err := model.UpdataUserNodePermissions(userId, nodeId, permission); err != nil {
		return err
	}

	// Update the root node of the user.
	if err := updateUserRootNode(userId); err != nil {
		return err
	}
	return nil
}

func UpdateGroupPermissionsToNode(groupId, nodeId interface{}, permission int) error {
	if err := model.UpdataGroupNodePermissions(userId, nodeId, permission); err != nil {
		return err
	}

	// The permission of group is changed, the user's root node of the group is maybe changed too.
	if err := updateUsersOfGroupRootNode(groupId); err != nil {
		return err
	}
	return nil
}

func updateUserRootNode(userId interface{}) error {

}

func updateUsersOfGroupRootNode(groupId interface{}) error {

}
