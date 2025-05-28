package models

type Product struct {
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	CategoryID   int           `json:"category_id,omitempty"`
	Price        float64       `json:"price,omitempty"`
	Description  string        `json:"description,omitempty"`
	DiscountName *string       `json:"discount_name,omitempty"` // pointer untuk handle NULL
	Discount     *float64      `json:"discount,omitempty"`      // pointer untuk handle NULL
	TotalOrder   int           `json:"total_order,omitempty"`
	Sizes        []ProductSize `json:"sizes,omitempty"`
	Images       []string      `json:"images,omitempty"`
	TotalRatings int           `json:"total_ratings,omitempty"`
	CategoryName string        `json:"category_name,omitempty"`
	IsDeleted    bool          `json:"isdeleted"`
}

type ProductSize struct {
	ID    int    `json:"id"`
	Name  string `json:"size"`
	Stock int    `json:"stock"`
}

// Products adalah slice dari Product
type Products []Product

type ProductQueryParams struct {
	Search   string   `form:"search"`
	Category []string `form:"category"`
	Discount *int     `form:"discount"`
	Options  string   `form:"options"`
	Min      int      `form:"min-price"`
	Max      int      `form:"max-price"`
	Page     int      `form:"page"`
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
