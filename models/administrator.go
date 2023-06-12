package models

import "gorm.io/gorm"

type Administrator struct {
	gorm.Model
	PhotoProfile     string `json:"photo_profile" form:"photo_profile" gorm:"default:'https://res.cloudinary.com/dvexlihfn/image/upload/v1686546113/go_bedu/mlc5oequ9xjvtm0w8kqb.jpg'"`
	Nama             string `json:"nama" form:"nama"`
	Email            string `json:"email" form:"email" validate:"required,email"`
	Username         string `json:"username" form:"username" validate:"required"`
	Password         string `json:"password" form:"password" validate:"required"`
	Role             string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null"`
	VerificationCode string
	OTP              int
	OTPReq           bool      `gorm:"not null"`
	Verified         bool      `gorm:"not null"`
	Token            string    `json:"-" gorm:"-"`
	Articles         []Article `json:"articles" form:"articles" gorm:"foreignKey:AdministratorID"`
}
