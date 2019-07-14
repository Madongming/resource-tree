package controller

type ResourceTree struct {
	Node   ResourceTreeNode `json:"node"`
	Childs []ResourceTree   `json:"childs"`
}

type ResourceTreeNode struct {
	ID          int64  `json:"id`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	ParentId    int64  `json:"parent_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	EnName      string `json:"enName`
	Key         string `json:"key"`
}

type NodeUserOrGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserPermission struct {
	ID        int   `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	User      IDRW  `json:"users"`
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
