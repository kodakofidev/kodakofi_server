package models

type UserReq struct {
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required,min=8"`
}
