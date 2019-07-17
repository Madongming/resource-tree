package model

import (
	. "global"
)

func CreateUser(name, cnName string) error {
	user := &User{
		Name:   name,
		CnName: cnName,
	}

	return user.create()
}

func CreateGroup(name, cnName string) error {
	group := &Group{
		Name:   name,
		CnName: cnName,
	}

	return group.create()
}

func AddUserToGroup(userId, groupId interface{}) error {
	ug := &UserGroup{
		UserID:  userId.(int),
		GroupID: groupId.(int),
	}
	return ug.create()
}

func AddUserToNode(userId, nodeId interface{}, permission int) error {
	un := &UserPermission{
		ReadWriteMask: permission,
		NodeID:        nodeId.(int),
		UserID:        userId.(int),
	}
	return un.create()
}

func AddGroupToNode(groupId, nodeId interface{}, permission int) error {
	gn := &UserPermission{
		ReadWriteMask: permission,
		NodeID:        nodeId.(int),
		GroupID:       groupId.(int),
	}
	return gn.create()
}

func (u *User) create() error {
	return DB().Create(u).Error
}

func (g *Group) create() error {
	return DB().Create(g).Error
}

func (ug *UserGroup) create() error {
	return DB().Create(ug).Error
}

func (gn *GroupPermission) create() error {
	return DB().Create(gn).Error
}

func (un *UserPermission) create() error {
	return DB().Create(un).Error
}
