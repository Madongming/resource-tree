package service

import (
	"model"
)

func UpdateUserPermissionsToNode(userId, nodeId interface{}, permission int) error {
	return model.UpdataUserNodePermissions(userId, nodeId, permission)
}

func UpdateGroupPermissionsToNode(groupId, nodeId interface{}, permission int) error {
	return model.UpdataGroupNodePermissions(userId, nodeId, permission)
}
