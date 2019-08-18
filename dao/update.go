package model

import (
	"cache"
	. "global"
)

// 更新节点名称，可选参数同于权限验证。如果想验证则传入userid
func UpdateNodeName(name string, nodeId interface{}, userId ...interface{}) error {
	// 更新缓存中的数据，将要在缓存中遍历节点树
	casResource()

	if userId != nil && len(userId) != 0 {
		i, err := isUserHasNodePermission(userId[0], nodeId)
		if err != nil {
			return err
		}
		if !i {
			return ERR_PERMISSION_DENY
		}
	}

	node, err := getNodeById(nodeId)
	if err != nil {
		return err
	}

	node.Name = name
	return node.Update()
}

func UpdataUserNodePermissions(userId, nodeId interface{}, permissions int) error {
	userPermission, err := getNodeUserPermission(userId, nodeId)
	if err != nil {
		return err
	} else if userPermission == nil {
		return nil
	}
	userPermission.ReadWriteMask = permissions
	return userPermission.update()
}

func UpdataGroupNodePermissions(groupId, nodeId interface{}, permissions int) error {
	groupPermission, err := getNodeGroupPermission(groupId, nodeId)
	if err != nil {
		return err
	} else if groupPermission == nil {
		return nil
	}
	groupPermission.ReadWriteMask = permissions
	return groupPermission.update()
}
