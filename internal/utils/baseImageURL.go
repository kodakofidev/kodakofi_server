package utils

import "os"

var BaseImgProfileURL string
var BaseImgProductURL string

func InitConfigProfile() {
	BaseImgProfileURL = os.Getenv("BASE_IMAGE_USER_URL")
	if BaseImgProfileURL == "" {
		BaseImgProfileURL = "http://localhost:8080/public/" // fallback local
	}
}

func InitConfigProduct() {
	BaseImgProductURL = os.Getenv("BASE_IMAGE_PRODUCT_URL")
	if BaseImgProductURL == "" {
		BaseImgProductURL = "http://localhost:8080/public/product-images/" // fallback local
	}
}
