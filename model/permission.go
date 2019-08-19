package model

import (
	"time"

	. "github.com/Madongming/resource-tree/global"
)

// DB recorder.
type DBUserPermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int  `gorm:"index"`
	UserID        int  `gorm:"index"`
}

// DB recorder.
type DBGroupPermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int  `gorm:"index"`
	GroupID       int  `gorm:"index"`
}

func (gn *DBGroupPermission) Create() error {
	return DB().Create(gn).Error
}

func (un *DBUserPermission) Create() error {
	return DB().Create(un).Error
}

func (gp *DBGroupPermission) Update() error {
	return DB().Save(gp).Error
}

func (up *DBUserPermission) Update() error {
	return DB().Save(up).Error
}

// About the ReadWriteMask
// 00 00 00 00 00 00 00 00 ...
// Every two bits is a group. Define option permissions.
// Etc: ...00 00 00 01<option 1> 11<option 2> means: <option 1> is readable <option 2> is readable&writable.
