package dao

import (
	. "github.com/Madongming/resource-tree/global"
)

// 更新节点名称，可选参数同于权限验证。如果想验证则传入userid
func UpdateNodeName(name, cnName string, nodeId interface{}, userId ...interface{}) error {
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
	node.CnName = cnName
	if err := node.Update(); err != nil {
		return err
	}
	casVersion()
	return nil
}

func UpdataUserNodePermissions(userId, nodeId interface{}, permissions int) error {
	userPermission, err := getNodeUserPermission(userId, nodeId)
	if err != nil {
		return err
	} else if userPermission == nil {
		return ERR_PERMISSION_DENY
	}
	userPermission.ReadWriteMask = uint(permissions)
	return userPermission.Update()
}

func UpdataGroupNodePermissions(groupId, nodeId interface{}, permissions int) error {
	groupPermission, err := getNodeGroupPermission(groupId, nodeId)
	if err != nil {
		return err
	} else if groupPermission == nil {
		return ERR_PERMISSION_DENY
	}
	groupPermission.ReadWriteMask = uint(permissions)
	return groupPermission.Update()
}
