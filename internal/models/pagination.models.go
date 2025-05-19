package models

type Pagination struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Links      map[string]string `json:"links,omitempty"`
}

type PaginatedResponse struct {
	Data       Products   `json:"data"`
	Pagination Pagination `json:"pagination"`
}
