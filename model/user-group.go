package model

import (
	"time"

	. "github.com/Madongming/resource-tree/global"
)

// DB recorder.
type DBUser struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(128)"`
	CnName    string `gorm:"type:varchar(378)"`
}

// DB recorder.
type DBGroup struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(128)"`
	CnName    string `gorm:"type:varchar(378)"`
}

// DB recorder.
type DBUserGroup struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int `gorm:"primary_key"`
	GroupID   int `gorm:"primary_key"`
}

func (u *DBUser) Create() error {
	return DB().Create(u).Error
}

func (g *DBGroup) Create() error {
	return DB().Create(g).Error
}

func (ug *DBUserGroup) Create() error {
	return DB().Create(ug).Error
}

func (u *DBUser) Update() error {
	return DB().Save(u).Error
}

func (g *DBGroup) Update() error {
	return DB().Save(g).Error
}

func (ug *DBUserGroup) Update() error {
	return DB().Save(ug).Error
}
