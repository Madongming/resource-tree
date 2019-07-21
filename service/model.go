package service

type Base struct {
	ID       int64  `json:"id"`
	RootNode int64  `json:"rootNode"`
	Name     string `json:"name"`
	CnName   string `json:"cnName"`
}

type User struct {
	Base
}

type Group struct {
	Base
}
