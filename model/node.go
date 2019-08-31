package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/Madongming/resource-tree/global"
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
	Name        string `gorm:"type:varchar(128);unique_index"`
	CnName      string `gorm:"type:varchar(378)"`
	Key         string `gorm:"type:varchar(512)"`
	Tags        string `gorm:"type:varchar(1024)"`
}

// In buffer
type ResourceNode struct {
	ID          int    `json:"id"`
	Parent      int    `json:"parent"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	CnName      string `json:"cnName"`
	Key         string `json:"key"`
	Tags        string `json:"tags"`
}

func (n *DBResourceNode) Create() error {
	return DB().Create(n).Error
}

func (n *DBResourceNode) Update() error {
	return DB().Save(n).Error
}

func (n *DBResourceNode) Delete() error {
	return DB().Delete(n).Error
}

// 其他属性做为非必选项
func (n *DBResourceNode) SetCnName(cnName interface{}, defaultName string) {
	cn, ok := cnName.(string)
	if ok && cnName != "" {
		n.CnName = cn
	} else {
		n.CnName = defaultName
	}
}

func (n *DBResourceNode) SetKey(key interface{}, defaultKey string) {
	k, ok := key.(string)
	if ok && k != "" {
		n.Key = k
	} else {
		n.Key = defaultKey
	}
}

func (n *DBResourceNode) SetTags(tags interface{}, defaultTags string) {
	t, ok := tags.(string)
	if ok && t != "" {
		n.Tags = t
	} else {
		n.Tags = defaultTags
	}
}

// 方便fmt.Print()打印
func (n *ResourceNode) String() string {
	b, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("%#v", n)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")

	if err != nil {
		return fmt.Sprintf("%#v", n)
	}

	return out.String()
}
