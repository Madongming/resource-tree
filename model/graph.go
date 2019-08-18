package model

import (
	"time"
)

// DB recorder.
type DBResourceRelationship struct {
	ID                   int `gorm:"primary_key"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
	SourceResourceNodeID int
	TargetResourceNodeID int
}

type ResourceEdge struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

type Graph struct {
	Nodes []*ResourceNode `json:"nodes"`
	Edges []*ResourceEdge `json:"edges"`
}
