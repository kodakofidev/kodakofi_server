package utils

import "os"

var BaseImgProfileURL string

func InitConfig() {
	BaseImgProfileURL = os.Getenv("BASE_IMAGE_URL")
	if BaseImgProfileURL == "" {
		BaseImgProfileURL = "http://localhost:8080/public/profile-images/" // fallback local
	}
}
