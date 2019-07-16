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
}
