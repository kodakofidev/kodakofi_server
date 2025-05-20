package models

import (
	"mime/multipart"
	"time"
)

type ProfileUser struct {
	Fullname string `json:"fullname" form:"fullname"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	// Email    string `json:"email" form:"email"`

	// Password string `json:"-" form:"password" binding:"required,min=8"`
}

type ProfileForm struct {
	ProfileUser
	ProfileImage *multipart.FileHeader `form:"profileImage"`
	// NewPass      string                `json:"-" form:"newPassword" binding:"required,min=8"`
	CurrentPassword string `json:"currentPassword" form:"currentPassword" binding:"required_with=NewPassword"`
	NewPassword     string `json:"newPassword" form:"newPassword" binding:"omitempty,min=8,max=72"`
}

type Profile struct {
	ProfileUser
	ProfileImage string    `json:"profileImage"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Profiles []Profile
