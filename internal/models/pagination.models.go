package models

type Pagination struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Link       map[string]string `json:"link"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}
