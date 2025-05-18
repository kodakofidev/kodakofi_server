package models

import "time"

type Order struct {
}

type OrderHistory struct {   
	TransactionCode string    `json:"transaction_code"`
	Date            time.Time `json:"date"`
	GrandTotal      int       `json:"grand_total"`
	OrderId         int       `json:"order_id"`
	Status          string    `json:"status"`
}

type OrderHistories []OrderHistory 

type Orders []Order
