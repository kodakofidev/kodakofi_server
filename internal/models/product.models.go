package models

type Product struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	CategoryID   int      `json:"category_id"`
	Price        float64  `json:"price"`
	Description  string   `json:"description"`
	DiscountName *string  `json:"discount_name,omitempty"` // pointer untuk handle NULL
	Discount     *float64 `json:"discount,omitempty"`      // pointer untuk handle NULL
	TotalOrder   int      `json:"total_order"`
	Images       []string `json:"images"`
	TotalRatings int      `json:"total_ratings"`
	CategoryName string   `json:"category_name"`
}

// Products adalah slice dari Product
type Products []Product

type ProductQueryParams struct {
	Page     int    `json:"page" form:"page" binding:"numeric"`
	Search   string `json:"search" form:"search"`
	Options  string `json:"options" form:"options"`
	Category string `json:"category" form:"category"`
	Discount string `json:"discount" form:"discount"`
	Min      int    `json:"min-price" form:"min-price" binding:"min=0"`
	Max      int    `json:"max-price" form:"max-price"`
}
