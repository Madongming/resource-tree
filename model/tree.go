package model

import (
	"time"
)

type ResourceTreeNode struct {
	ID          int `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Parent      int
	Description string `gorm:"type:varchar(1024)"`
	Level       int    // 0 root; 1 child; 2 child...
	Name        string `gorm:"type:varchar(378)"`
	EnName      string `gorm:"type:varchar(128)"`
}

type UserPermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int
	UserID        int
}

type GroupPermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int
	GroupID       int
}

type User struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RootNode  int
	Name      string `gorm:"type:varchar(378)"`
	EnName    string `gorm:"type:varchar(128)"`
}

type Group struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RootNode  int
	Name      string `gorm:"type:varchar(378)"`
	EnName    string `gorm:"type:varchar(128)"`
}

type Tag struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeId    int
	Text      `gorm:"type:varchar(4096)"`
}

// About the ReadWriteMask
// 00 00 00 00 00 00 00 00 ...
// Every two bits is a group. Define option permissions.
// Etc: ...00 00 00 01<option 1> 11<option 2> means: <option 1> is readable <option 2> is readable&writable.
