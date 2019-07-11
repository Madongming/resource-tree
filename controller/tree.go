package controller

type ResourceTreeNode struct {
	ID          int64              `json:"id`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
	Childs      []ResourceTreeNode `json"childs"`
	Description string             `jsonL"description"`
	Level       int                `json:"level"`
	Name        string             `json:"name"`
	EnName      string             `json:"enName`
}

type UserPermission struct {
	ID        int   `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	Users     IDRW  `json:"users"`
}

type GroupPermission struct {
	ID        int   `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	Groups    IDRW  `json:"groups"`
}

type IDRW struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Readable bool   `json:"readable"`
	Writable bool   `json:"writable"`
}
