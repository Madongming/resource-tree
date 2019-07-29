package service

import (
	"cache"
	"model"
)

func GetAllTree() (*cache.ResourceTree, error) {
	return cache.Cache.Tree, nil
}

// 获取user以及其所在组的所有可读节点，用map当set来使用
func GetUserNodes(userId int) (map[int]struct{}, error) {
	groups, err := model.GetUserGroups(userId)
	if err != nil {
		return nil, err
	}

	groupNodeIds, err := model.GetGroupsNodeIds(groups)
	if err != nil {
		return nil, err
	}
	userNodeIds, err := model.GetUserNodes(userId)
	if err != nil {
		return nil, err
	}

	// 最多为两者之和
	resultSet := make(map[int]struct{}, len(groupNodeIds)+len(userNodeIds))

	for i := range groupNodeIds {
		resultSet[groupNodeIds[i]] = struct{}{}
	}
	for i := range userNodeIds {
		resultSet[userNodeIds[i]] = struct{}{}
	}
	return resultSet, nil
}

func GetTree(userId interface{}, isFull ...bool) (*cache.ResourceTree, error) {
	isCas, err := cache.CasResource()
	if err != nil {
		return nil, err
	}
	if isFull == nil || len(isFull) == 0 || isFull[0] {
		return GetAllTree()
	}
	if !isCas {
		// 从缓存中取
	}

	// 如果缓存还存在，删除缓存，重建用户树
	tree := copyUserTree()

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
