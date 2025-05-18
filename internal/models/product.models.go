package models

type Product struct {
	ID           string   `json:"id,omitempty" form:"id" db:"id"`
	Name         string   `json:"name,omitempty" form:"name" db:"name"`
	CategoryID   int      `json:"category_id,omitempty" form:"category_id"  db:"category_id"`
	Price        int      `json:"price,omitempty" form:"price"  db:"price"`
	Description  string   `json:"description,omitempty" form:"description"  db:"description"`
	DiscountName *string  `json:"discount_name,omitempty" form:"discount_name"  db:"discount_name"`
	Discount     *int     `json:"discount,omitempty" form:"discount"  db:"discount"`
	TotalOrder   *int     `json:"total_order,omitempty" form:"total_order"  db:"total_order"`
	Images       []string `json:"images,omitempty" form:"images"  db:"images"`
	TotalRatings int      `json:"total_ratings,omitempty" form:"total_ratings"  db:"total_ratings"`
	CategoryName string   `json:"category_name,omitempty" form:"category_name"  db:"category_name"`
}

type Products []Product

type ProductQueryParams struct {
	Page     int    `json:"page" form:"page" binding:"numeric"`
	Search   string `json:"search" form:"search"`
	Name     string `json:"name" form:"name"`
	Options  string `json:"options" form:"options"`
	Category string `json:"category" form:"category"`
	Discount string `json:"discount" form:"discount"`
	Min      int    `json:"min-price" form:"min-price" binding:"min=0"`
	Max      int    `json:"max-price" form:"max-price"`
}
