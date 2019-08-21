package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/Madongming/resource-tree/global"
)

// DB recorder.
type DBResourceRelationship struct {
	ID                   int `gorm:"primary_key"`
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

func (g *Graph) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%#v", g)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")

	if err != nil {
		return fmt.Sprintf("%#v", g)
	}

	return out.String()
}
