package model

import (
	"time"
)

// DB recorder.
type DBResourceNode struct {
	ID          int `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Parent      int
	Description string `gorm:"type:varchar(1024)"`
	Level       int    // 0 root; 1 child; 2 child...
	Name        string `gorm:"type:varchar(128)"`
	CnName      string `gorm:"type:varchar(378)"`
	Key         string `gorm:"type:varchar(512)"`
}

// In buffer
type ResourceNode struct {
	ID          int    `json:"id"`
	Parent      int    `json:"parent"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Name        string `json:"Name"`
	CnName      string `json:"cnName"`
	Key         string `json:"key"`
}
