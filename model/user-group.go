package model

import (
	"time"
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
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	GroupID   int
}
