package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username         string `json:"username" form:"username" validate:"required"`
	Password         string `json:"password" form:"password"`
	FullName         string `json:"fullname" form:"fullname"`
	Email            string `json:"email" form:"email"`
	Role             string `json:"role" form:"role" gorm:"type:enum('User');default:'User'; not-null"`
	VerificationCode string
	OTP              int
	OTPReq           bool   `gorm:"not null"`
	Verified         bool   `gorm:"not null"`
	Token            string `json:"-" gorm:"-"`
	PhotoProfile     string `json:"photo_profile" form:"photo_profile" gorm:"default:'https://res.cloudinary.com/dvexlihfn/image/upload/v1686546113/go_bedu/mlc5oequ9xjvtm0w8kqb.jpg'"`
}
