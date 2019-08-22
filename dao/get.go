package dao

import (
	"strings"

	"github.com/Madongming/resource-tree/cache"
	. "github.com/Madongming/resource-tree/global"
	"github.com/Madongming/resource-tree/model"
)

// 获取组内的所有用户
func GetGroupUsers(groupId interface{}) ([]*model.DBUser, error) {
	var userGroups []*model.DBUserGroup
	if result := DB().
		Where("`group_id` = ?", groupId).
		Find(&userGroups); result.Error != nil {
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
func GetTreeNodeGroups(nodeId interface{}) ([]*model.DBGroup, error) {
	var permissions []*model.DBGroupPermission
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

// 获取指定用户的树，可选获取不分权限的整棵树。
func GetTree(userId interface{}, isFull ...bool) (*model.Tree, error) {
	isCas, err := casResource()
	if err != nil {
		return nil, err
	}
	if isFull != nil && len(isFull) != 0 && isFull[0] {
		// 不做权限校验，获取整棵树
		return getAllTree()
	}
	if !isCas {
		// 从缓存中取
		tree, err := cache.UserTreeLRU.Get(userId.(int))
		if err != nil {
			if err.Error() == ERR_CACHE_KEY_NOT_EXIST.Error() {
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

// 获取指定节点的图
func GetNodeGraph(nodeId interface{}) (*model.Graph, error) {
	version, err := getCurrentEdgeVersion()
	if err != nil {
		return nil, err
	}

	if value, err := cache.ResourceLRU.Get(nodeId.(int)); err == nil {
		if version == cache.ResourceLRU.Version {
			// 从缓存中取
			return value.(*model.Graph), nil
		}
	}

	keys, nodes, rrs, err := getUserGraph(nodeId)
	if err != nil {
		return nil, err
	}
	graph := &model.Graph{
		Nodes: nodes,
		Edges: rrs,
	}
	cache.ResourceLRU.Version = version
	cache.ResourceLRU.Set(keys, graph)
	return graph, nil
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

// 获取所有节点列表
func GetAllResourceNodes(userId interface{}, isFull ...bool) ([]*model.ResourceNode, error) {
	tree, err := GetTree(userId, isFull...)
	if err != nil {
		return nil, err
	}
	results := make([]*model.ResourceNode, cache.Tree.Size)
	var index int
	levelOrderTraverse(tree, results, cache.Tree.Size, &index)
	return results[0:index], nil
}
