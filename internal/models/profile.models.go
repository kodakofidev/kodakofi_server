package models

import "mime/multipart"

type ProfileUser struct {
	Id       string `json:"id:"`
	Fullname string `json:"fullname" form:"fullname"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
}

type ProfileForm struct {
	ProfileUser
	ProfileImage *multipart.FileHeader `form:"profileImage"`
}

type Profile struct {
	ProfileUser
	ProfileImage string `json:"profileImage"`
}

type Profiles []Profile
