package model

type Tree struct {
	Node   *ResourceNode `json:"node"`
	Childs []*Tree       `json:"childs"`
}
