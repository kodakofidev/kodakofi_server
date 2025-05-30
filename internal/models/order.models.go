package models

import "time"

type OrderItem struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
	SizeID    int    `json:"size_id"`
	IsIced    bool   `json:"is_iced"`
}

type OrderHistory struct {
	TransactionCode string    `json:"transaction_code"`
	Date            time.Time `json:"date"`
	GrandTotal      int       `json:"grand_total"`
	OrderId         int       `json:"order_id"`
	Status          string    `json:"status"`
	Path            *string   `json:"path"`
}

// type OrderHistories []OrderHistory

type CreateOrderRequest struct {
	Email            string      `json:"email"`
	Fullname         string      `json:"fullname"`
	Address          string      `json:"address"`
	DeliveryMethodID int         `json:"delivery_method_id"`
	PaymentMethodID  int         `json:"payment_method_id"`
	Items            []OrderItem `json:"items"`
}

type OrderItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Qty         int    `json:"qty"`
	Size        string `json:"size"`
	IsIced      bool   `json:"is_iced"`
}

type CreateOrderResponse struct {
	Email          string              `json:"email"`
	Fullname       string              `json:"fullname"`
	Address        string              `json:"address"`
	DeliveryMethod int                 `json:"delivery_method"`
	PaymentMethod  int                 `json:"payment_method"`
	Items          []OrderItemResponse `json:"items"`
	DeliveryFee    int                 `json:"delivery_fee"`
	Total          int                 `json:"total"`
	Tax            int                 `json:"tax"`
	TotalAmount    int                 `json:"total_amount"`
}

type TotalIncomeItemResponse struct {
	ProductName   string `json:"product_name"`    // diambil dari column product_id di table products_orders dan mereferensikan ke tabel products
	TotalItemSold int    `json:"total_item_sold"` // diambil dari column qty products_orders yaitu total seluruh product terjual pada rentang waktu tetentu yang dipilih (referensi created_at di tabel orders)
	Income        string `json:"income"`          // diambil dari column sub_total di tabel products_orders bedasarkan nama product yang sama
}

type SalesDataStatus struct {
	Pending    int `json:"pending"`    // jumlah status Pending dari column status di tabel orders
	Processing int `json:"processing"` // jumlah status Processing dari column status di tabel orders
	Completed  int `json:"completed"`  // jumlah status Completed dari column status di tabel orders
} // status diambil dari status per order dari tabel orders

type DailySoldItems struct {
	Date         string `json:"date"`          // diambil dari column created_at di tabel orders
	ProductsSold int    `json:"products_sold"` // julah seluruh product terjual di satu hari
}

type ProductSalesDataRes struct {
	Status            SalesDataStatus           `json:"status"`           // diambil dari column status_id di tabel orders
	TotalSoldItems    int                       `json:"total_sold_items"` // diambil dari column qty di products_orders dengan referensi tanggal yang dipilih dengan rentang tertentu pada tabel orders
	DailySoldItems    []DailySoldItems          `json:"daily_sold_items"` // data data jumlah produk terjual per hari dari dalam rentang waktu tanggal yang dipulih
	IncomeDataPerItem []TotalIncomeItemResponse `json:"income_data"`      // default langsung diurutkan dari item yang paling banyak terjual
	TotalData         int                       `json:"total_data"`       // total data dari IncomeDataPerItem, jadi hanya data IncomeDataPerItem yang berubah akibat data yang berlebihan
}

type OrderStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type UpdateOrderStatusReq struct {
	OrderID  int `json:"order_id"`
	StatusID int `json:"status_id"`
}

type UpdateOrderStatusRes struct {
	TransactionCode string    `json:"transaction_code"`
	OrderID         int       `json:"order_id"`
	Status          string    `json:"status"`
	UpdateAt        time.Time `json:"updated_at"`
}

type AdminOrderHistory struct {
	TransactionCode string              `json:"transaction_code"`
	CreatedAt       string              `json:"created_at"`
	Status          string              `json:"status"`
	TotalAmount     int                 `json:"total_amount"`
	OrderItems      []AdminOrderItemRes `json:"order_items"`
}

type AdminOrderHistories []AdminOrderHistory

type AdminOrderItemRes struct {
	ProductName string `json:"product_name"`
	Size        string `json:"size"`
	Temperature string `json:"temperature"`
	Qty         string `json:"qty"`
}

type DetailOrderRes struct {
	TransactionCode string `json:"transaction_code"` // diambil dari tabel transactions
	Fullname        string `json:"fullname"`         // diambil dari tabel orders
	Address         string `json:"address"`          // diambil dari tabel orders
	Phone           string `json:"phone"`            // diambil dari tabel profiles dengan referensid dari tabel users dan tabel orders yaitu user_id
	PaymentMethod   string `json:"payment_method"`   // diambil dari tabel payment_methods yang direfernsikan dari tabel orders dari column payment_methods_id
	Shipping        string `json:"shipping"`         // diambil dari tabel delivery_methods yang direfernsikan dari tabel orders dari column delivery_methods_id
	DeliveryFee     int    `json:"delivery_fee"`     // diambil dari column delivery_fee di tabel transactions
	Tax             int    `json:"tax"`              // diambil dari column tax pada tabel transactions
	Status          string `json:"status"`           // diambil dari tabel orders di column status_id dan merefernsikan ke tabel status
	TotalAmount     int    `json:"total_amount"`     // diambil dari column total_amount di tabel transactions
	CreatedAt       string `json:"created_at"`
	OrderItems      []ItemDetailOrder
}

type ItemDetailOrder struct {
	ProductID   string `json:"product_id"`   // diambil dari tabel products_orders
	ProductName string `json:"product_name"` // diambil dari product name dari referensi id pada product
	ProductImg  string `json:"product_img"`
	Qty         int    `json:"qty"`         // diambil dari tabel products_orders
	Size        string `json:"size"`        // diambil dari tabel products_orders
	Temperature string `json:"temperature"` // diambil dari tabel products_orders dari column is_iced dan dikonversi menjadi string di response
	SubTotal    int    `json:"sub_total"`   // diambil dari column sub_total dari tabel products_orders
}
