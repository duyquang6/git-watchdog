package dto

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"string"`
}

type PaginationMeta struct {
	Meta
	Total uint `json:"uint"`
}
