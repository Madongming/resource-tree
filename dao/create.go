package dao

import (
	. "github.com/Madongming/resource-tree/global"
	"github.com/Madongming/resource-tree/model"
	"github.com/Madongming/resource-tree/tools"
)

// 将一个节点加入树
func CreateNode(name, description string, userId, parentId int, opts ...interface{}) error {
	var (
		userPermission = 1
		groupPermisson = 1
		groupId        = 0
	)

	node := new(model.DBResourceNode)

	// opts 0: cn_name, 1: key, 2: user permission, 3: group id, 4, group permission.
	if opts == nil || len(opts) == 0 {
		node.CnName = name
		node.Key = tools.GetUuid()
	} else {
		switch len(opts) {
		case 1:
			node.CnName = opts[0].(string)
			node.Key = tools.GetUuid()
		case 2:
			node.CnName = opts[0].(string)
			node.Key = opts[1].(string)
		case 3:
			node.CnName = opts[0].(string)
			node.Key = opts[1].(string)
			userPermission = opts[2].(int)
		case 4:
			node.CnName = opts[0].(string)
			node.Key = opts[1].(string)
			userPermission = opts[2].(int)
			groupId = opts[3].(int)
		default:
			node.CnName = opts[0].(string)
			node.Key = opts[1].(string)
			userPermission = opts[2].(int)
			groupId = opts[3].(int)
			groupPermisson = opts[4].(int)
		}
	}

	if parentId == 0 {
		node.Level = 1
	} else {
		parentNode := new(model.DBResourceNode)
		if err := DB().First(parentNode, parentId).Error; err != nil {
			return err
		}
		node.Level = parentNode.Level + 1
	}
	if err := node.Create(); err != nil {
		return err
	}

	// Node版本加一
	casVersion()

	if groupId != 0 {
		if err := addGroupToNode(groupId, node.ID, groupPermisson); err != nil {
			return err
		}
	}
	return addUserToNode(userId, node.ID, userPermission)
}

// Add user when user first login.
func AddUser(name, cnName string) error {
	user, err := getUserByName(name)
	if user != nil {
		return nil
	} else if err != nil {
		return err
	}
	return createUser(name, cnName)
}

// Create a group.
func CreateGroup(name, cnName string) error {
	return createGroup(name, cnName)
}

// 将一个用户加入组
func AddUserToGroup(userId, groupId interface{}) error {
	return addUserToGroup(userId, groupId)
}

// 将一个节点授权给一个用户
func AddUserToNode(userId interface{}, nodeId int, permissions ...int) error {
	// Default permission is 1.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	return addUserToNode(userId, nodeId, permission)
}

// 将一个节点授权给组
func AddGroupToNode(groupId interface{}, nodeId int, permissions ...int) error {
	// Default permission is 1.
	permission := 1
	if permissions != nil && len(permissions) != 0 {
		permission = permissions[0]
	}

	return addGroupToNode(groupId, nodeId, permission)
}

// 加入一条资源的管理关系，可选参数userId如果传入，则判断两个节点是否都有权限
func AddNodeToNode(srcNodeId, tarNodeId interface{}, userId ...int) error {
	rr := new(model.DBResourceRelationship)
	rr.SourceResourceNodeID = srcNodeId.(int)
	rr.TargetResourceNodeID = tarNodeId.(int)
	if userId == nil || len(userId) == 0 {
		return rr.Create()
	}

	srci, err := isUserHasNodePermission(userId[0], srcNodeId.(int))
	if err != nil {
		return err
	}
	tari, err := isUserHasNodePermission(userId[0], tarNodeId.(int))
	if err != nil {
		return err
	}
	if srci && tari {
		return rr.Create()
	}
	casEdgeVersion()
	return ERR_PERMISSION_DENY
}
