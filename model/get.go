package model

import (
	"strings"

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
			return []*ResourceTreeNode{}, nil
		}
	}
	return resourceTreeNodes, nil
}

func GetUserById(id interface{}) (*User, error) {
	user := &User{}
	if result := DB().
		First(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func GetUsersByIds(ids []int) ([]*User, error) {
	var users []*User
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}
	idsIter := make([]interface{}, len(ids))
	idsQueryStr := make([]string, len(ids))
	for i := range ids {
		idsIter[i] = ids[i]
		idsQueryStr[i] = "?"
	}

	if result := DB().
		Where("`id` in ("+
			strings.Join(idsQueryStr, ",")+
			")",
			idsIter...).
		Find(&users); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func GetGroupById(id interface{}) (*Group, error) {
	group := &Group{}
	if result := DB().
		First(group); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return group, nil
}

func GetGroupsByIds(ids []int) ([]*Group, error) {
	var groups []*Group
	if ids == nil || len(ids) == 0 {
		return nil, nil
	}
	idsIter := make([]interface{}, len(ids))
	idsQueryStr := make([]string, len(ids))
	for i := range ids {
		idsIter[i] = ids[i]
		idsQueryStr[i] = "?"
	}

	if result := DB().
		Where("`id` in ("+
			strings.Join(idsQueryStr, ",")+
			")",
			idsIter...).
		Find(&groups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groups, nil
}

func GetUserByName(name string) (*User, error) {
	user := &User{}
	if result := DB().
		Where("`name` = ?", name).
		First(user); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return user, nil
}

func GetGroupUsers(groupId interface{}) ([]*User, error) {
	var userGroups []*UserGroup
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(userGroups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}

	userIds := make([]interface{}, len(userGroups))
	userQuery := make([]string, len(userGroups))
	for i := range userGroups {
		userIds[i] = userGroups[i].UserID
		userQuery[i] = "?"
	}

	var users []*User
	if result := DB().
		Where("`id` in ("+
			strings.Join(userQuery, ",")+
			")", userIds...).
		Find(&users); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func GetUserGroups(userId interface{}) ([]*Group, error) {
	var userGroups []*UserGroup
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(userGroups); result.Error != nil {
		if result.RecordNotFound() {
			return []*Group{}, nil
		}
		return nil, result.Error
	}

	groupIds := make([]interface{}, len(userGroups))
	groupQuery := make([]string, len(userGroups))
	for i := range userGroups {
		groupIds[i] = userGroups[i].UserID
		groupQuery[i] = "?"
	}

	var groups []*Group
	if result := DB().
		Where("`id` in ("+
			strings.Join(groupQuery, ",")+
			")", groupIds...).
		Find(&groups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groups, nil
}

func GetNodeUserPermission(userId, nodeId interface{}) (*UserPermission, error) {
	userPermission := &UserPermission{}
	if result := DB().
		Where("`user_id` = ? AND "+
			"`node_id` = ?",
			userId, nodeId).
		First(userPermission); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return userPermission, nil
}

func GetNodeGroupPermission(groupId, nodeId interface{}) (*GroupPermission, error) {
	groupPermission := &GroupPermission{}
	if result := DB().
		Where("`group_id` = ? AND "+
			"`node_id` = ?",
			groupId, nodeId).
		First(groupPermission); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	return groupPermission, nil
}

func GetTreeNodeUsers(nodeId interface{}) ([]*User, error) {
	var permissions []*UserPermission
	if result := DB().
		Where("`node_id` = ?", nodeId).
		Find(&permissions); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	userIds := make([]int, len(permissions))
	for i := range permissions {
		userIds[i] = permissions[i].UserID
	}

	return GetUsersByIds(userIds)
}

func GetTreeNodeGroups(groupId interface{}) ([]*User, error) {
	var permissions []*GroupPermission
	if result := DB().
		Where("`node_id` = ?", nodeId).
		Find(&permissions); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, result.Error
	}
	groupIds := make([]int, len(permissions))
	for i := range permissions {
		groupIds[i] = permissions[i].GroupID
	}

	return GetGroupsByIds(groupIds)
}

func GetGroupNodeIds(groupId interface{}) ([]int, error) {
	var groupPermissions []*GroupPermission
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(&groupPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, result.Error
	}
	nodeIds := make([]int, len(groupPermissions))
	for i := range groupPermissions {
		nodeIds[i] = groupPermissions[i].NodeID
	}
	return nodeIds
}

func GetUserNodeIds(userId interface{}) ([]int, error) {
	var userPermissions []*UserPermission
	if result := DB().
		Where("`user_id` = ?", userId).
		Find(&userPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, result.Error
	}
	nodeIds := make([]int, len(userPermissions))
	for i := range userPermissions {
		nodeIds[i] = userPermissions[i].NodeID
	}
	return nodeIds
}

func GetGroupsNodeIds(groups []*Group) ([]int, error) {
	if groups == nil || len(groups) == 0 {
		return []int{}, nil
	}
	groupIds := make([]interface{}, len(groups))
	queryArray := make([]string, len(groups))
	for i := range groups {
		groupIds[i] = groups[i].ID
		queryArray[i] = "?"
	}

	var groupPermissions []*GroupPermission
	if result := DB().
		Where("`group_id` in (" +
			strings.Join(queryArray, ","+
				")", groupIds...)).
		Find(&groupPermissions); result.Error != nil {
		if result.RecordNotFound() {
			return []int{}, nil
		}
		return nil, err
	}
	nodeIds := make([]int, len(groupPermissions))
	for i := range groupPermissions {
		nodeIds[i] = groupPermissions[i].NodeID
	}
	return nodeIds, nil
}
