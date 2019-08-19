package model

import (
	"time"

	. "github.com/Madongming/resource-tree/global"
)

// DB recorder.
type DBResourceRelationship struct {
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
	SourceResourceNodeID int `gorm:"index"`
	TargetResourceNodeID int `gorm:"index"`
}

type ResourceEdge struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

type Graph struct {
	Nodes []*ResourceNode `json:"nodes"`
	Edges []*ResourceEdge `json:"edges"`
}

func (rr *DBResourceRelationship) Create() error {
	return DB().Create(rr).Error
}

func (rr *DBResourceRelationship) Delete() error {
	return DB().Delete(rr).Error
}
