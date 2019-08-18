package model

import (
	"time"
)

// DB recorder.
type DBUserTreePermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int
	UserID        int
}

// DB recorder.
type DBGroupPermission struct {
	ID            int `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReadWriteMask uint // See Last
	NodeID        int
	GroupID       int
}

// About the ReadWriteMask
// 00 00 00 00 00 00 00 00 ...
// Every two bits is a group. Define option permissions.
// Etc: ...00 00 00 01<option 1> 11<option 2> means: <option 1> is readable <option 2> is readable&writable.
