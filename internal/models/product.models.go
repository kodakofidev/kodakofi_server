package models

type Product struct {
	ID            int `json:"page" form:"page" binding:"numeric" db:`
	Name          string
	CategoryID    int
	Price         int
	Description   string
	DiscountName  string
	Discount      int
	TotalOrder    int
	Images        []string
	TotalRatings  int
	AverageRating int
	CategoryName  string
}

type Products []Product

type ProductQueryParams struct {
	Page     int    `json:"page" form:"page" binding:"numeric"`
	Search   string `json:"search" form:"search"`
	Name     string `json:"name" form:"name"`
	Options  string `json:"options" form:"options"`
	Category string `json:"category" form:"category"`
	Discount string `json:"discount" from:"discount"`
	Min      int    `json:"min-price" from:"min-price"`
	Max      int    `json:"max-price" from:"max-price"`
}
