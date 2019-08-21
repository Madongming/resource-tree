package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Tree struct {
	Node   *ResourceNode `json:"node"`
	Childs []*Tree       `json:"childs"`
}

// 方便fmt.Print()打印
func (t *Tree) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%#v", t)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")

	if err != nil {
		return fmt.Sprintf("%#v", t)
	}

	return out.String()
}
