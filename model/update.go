package model

import (
	. "global"
)

func UpdateUserRootNode(userId, rootNode interface{}) error {
	user, err := GetUserById(userId)
	if err != nil {
		return err
	} else if user == nil {
		return nil
	}
	user.RootNode = rootNode.(int)
	return user.update()
}

func UpdataUserNodePermissions(userId, nodeId interface{}, permissions int) error {
	userPermission, err := GetNodeUserPermission(userId, nodeId)
	if err != nil {
		return err
	} else if userPermission == nil {
		return nil
	}
	userPermission.ReadWriteMask = permissions
	return userPermission.update()
}

func UpdataGroupNodePermissions(groupId, nodeId interface{}, permissions int) error {
	groupPermission, err := GetNodeGroupPermission(groupId, nodeId)
	if err != nil {
		return err
	} else if groupPermission == nil {
		return nil
	}
	groupPermission.ReadWriteMask = permissions
	return groupPermission.update()
}

func (u *User) update() error {
	return DB().Save(u).Error
}

func (g *Group) update() error {
	return DB().Save(g).Error
}

func (ug *UserGroup) update() error {
	return DB().Save(ug).Error
}

func (gp *GroupPermission) update() error {
	return DB().Save(gp).Error
}

func (up *UserPermission) update() error {
	return DB().Save(up).Error
}
