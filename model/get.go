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

func GetUserGroups(userId interface{}) ([]*User, error) {
	var userGroups []*UserGroup
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(userGroups); result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
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
