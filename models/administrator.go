package models

import "gorm.io/gorm"

type Administrator struct {
	gorm.Model
	Nama             string `json:"nama" form:"nama"`
	Email            string `json:"email" form:"email" validate:"required,email"`
	Password         string `json:"password" form:"password" validate:"required"`
	Role             string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null"`
	VerificationCode string
	Verified         bool      `gorm:"not null"`
	Token            string    `json:"-" gorm:"-"`
	Articles         []Article `json:"articles" form:"articles" gorm:"foreignKey:AdministratorID"`
}
