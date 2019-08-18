package dao

import (
	"strings"

	"cache"
	. "global"
	"model"
)

// 获取组内的所有用户
func GetGroupUsers(groupId interface{}) ([]*model.DBUser, error) {
	var userGroups []*model.DBUserGroup
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

	var users []*model.DBUser
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

// 获取节点有权限的用户，只有是显示设定，而不是继承来的才有显示。
func GetTreeNodeUsers(nodeId interface{}) ([]*model.DBUser, error) {
	var permissions []*model.DBUserPermission
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

	return getUsersByIds(userIds)
}

// 获取节点有权限的组，只有是显示设定，而不是继承来的才有显示。
func GetTreeNodeGroups(groupId interface{}) ([]*model.DBGroup, error) {
	var permissions []*model.GroupPermission
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

	return getGroupsByIds(groupIds)
}

// 获取指定用户的树，可选获取不分权限的树。
func GetTree(userId interface{}, isFull ...bool) (*model.Tree, error) {
	isCas, err := casResource()
	if err != nil {
		return nil, err
	}
	if isFull == nil || len(isFull) == 0 || isFull[0] {
		// 不做权限校验，获取整棵树
		return getAllTree()
	}
	if !isCas {
		// 从缓存中取
		tree, err := cache.UserTreeLRU.Get(userId.(int))
		if err != nil {
			if err == ERR_CACHE_KEY_NOT_EXIST {
				goto SET_CACHE_AND_RETURN
			}
			return nil, err
		}
		return tree.(*model.Tree), nil
	}

SET_CACHE_AND_RETURN:
	// 没有缓存
	tree, err := makeUserTreeIndex(userId)
	if err != nil {
		return nil, err
	}
	cache.UserTreeLRU.Set(userId.(int), tree)
	return tree, nil

}

// 获取一个用户对一个节点的权限。
func GetUserPermission(userId, nodeId int) (int64, error) {
	permission, err := getNodeUserPermission(userId, nodeId)
	if err != nil {
		return 0, err
	}

	return int64(permission.ReadWriteMask), nil
}

// 获取一个组对一个节点的权限
func GetGroupPermission(groupId, nodeId int) (int64, error) {
	permission, err := getNodeGroupPermission(groupId, nodeId)
	if err != nil {
		return int64(0), err
	}

	return int64(permission.ReadWriteMask), nil
}
