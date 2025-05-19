package models

type Product struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	CategoryID   int           `json:"category_id"`
	Price        float64       `json:"price"`
	Description  string        `json:"description"`
	DiscountName *string       `json:"discount_name,omitempty"` // pointer untuk handle NULL
	Discount     *float64      `json:"discount,omitempty"`      // pointer untuk handle NULL
	TotalOrder   int           `json:"total_order"`
	Sizes        []ProductSize `json:"sizes"`
	Images       []string      `json:"images"`
	TotalRatings int           `json:"total_ratings"`
	CategoryName string        `json:"category_name"`
}

type ProductSize struct {
	ID    int    `json:"id"`
	Name  string `json:"size"`
	Stock int    `json:"stock"`
}

// Products adalah slice dari Product
type Products []Product

type ProductQueryParams struct {
	Search   string `form:"search"`
	Category string `form:"category"`
	Discount string `form:"discount"`
	Options  string `form:"options"`
	Min      int    `form:"min"`
	Max      int    `form:"max"`
	Page     int    `form:"page"`
}

type ProductRequest struct {
	Id          string   `json:"id" form:"id"  db:"id"`
	Name        *string  `json:"name" form:"name"`
	Price       *int     `json:"price" form:"price"`
	Description *string  `json:"description" form:"description"`
	Stock       *int     `json:"stock" form:"stock" `
	CategoryID  *int     `json:"category_id" form:"category_id"`
	Size        []int    `json:"size" form:"size"`
	KeepImages  []string `json:"keep_images" form:"keep_images" `
}
