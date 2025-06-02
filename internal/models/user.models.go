package models

import "mime/multipart"

type User struct {
	ID         string `json:"id"`
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

type UserDetailsRes struct {
	User
	UpdatedAt string `json:"updated_at"`
}

type UpdateUserByAdminReq struct {
	ID       string                `form:"-"`
	Fullname string                `form:"fullname"`
	Phone    string                `form:"phone"`
	Role     string                `form:"role"`
	Address  string                `form:"address"`
	Image    *multipart.FileHeader `form:"image"`
}
