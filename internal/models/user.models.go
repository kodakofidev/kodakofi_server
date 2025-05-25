package models

type User struct {
	Fullname   string `json:"fullname"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Image      string `json:"image"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

type Users []User
