package models

import "time"

type UserReq struct {
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required,min=8"`
	Fullname string `json:"fullname" db:"fullname" binding:"omitempty"`
}

type UserRes struct {
	Email      string `json:"email" db:"email"`
	Role       string `json:"role" db:"role"`
	AuthID     string `json:"id" db:"auth_id"`
	Pass       string `json:"password" db:"password"`
	IsVerified bool   `json:"is_verified" db:"is_verified"`
	Fullname   string `json:"fullname" db:"fullname"`
}

type OTPVerificationReq struct {
	Email  string `json:"email" binding:"required"`
	OTP    string `json:"otp" binding:"required"`
	TypeID int    `json:"type_id" binding:"required"` // Added TypeID for flexible verification
}

type OTPData struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	TypeID    int       `json:"type_id" db:"type_id"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
}

type OTPSend struct {
	Email  string `json:"email" binding:"required"`
	TypeID int    `json:"type_id" binding:"required"`
}
