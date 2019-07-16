package model

import (
	. "global"
)

func (u *User) Save() error {
	return DB().Save(u).Error
}

func (g *Group) Save() error {
	return DB().Save(g).Error
}
