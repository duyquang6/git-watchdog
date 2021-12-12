package dto

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type PaginationMeta struct {
	Meta
	Total uint `json:"total"`
}
