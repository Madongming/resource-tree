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

func AddUserToNode(userId, nodeId interface{}, permissions int) error {
	un := &UserPermission{
		ReadWriteMask: permissions,
		NodeID:        nodeId.(int),
		UserID:        userId.(int),
	}
	return un.create()
}

func AddGroupToNode(groupId, nodeId interface{}, permissions int) error {
	gn := &UserPermission{
		ReadWriteMask: permissions,
		NodeID:        nodeId.(int),
		GroupID:       groupId.(int),
	}
	return gn.create()
}

func (u *User) create() error {
	return DB().Create(u).Error
}

func (u *User) update() error {
	return DB().Save(u).Error
}

func (g *Group) create() error {
	return DB().Create(g).Error
}

func (g *Group) update() error {
	return DB().Save(g).Error
}

func (ug *UserGroup) create() error {
	return DB().Create(ug).Error
}

func (ug *UserGroup) update() error {
	return DB().Save(ug).Error
}

func (gn *GroupPermission) create() error {
	return DB().Create(gn).Error
}

func (gn *GroupPermission) update() error {
	return DB().Save(gn).Error
}

func (un *UserPermission) create() error {
	return DB().Create(un).Error
}

func (un *UserPermission) update() error {
	return DB().Save(un).Error
}
