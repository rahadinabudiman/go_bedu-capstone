package models

import "github.com/jinzhu/gorm"

type Article struct {
	gorm.Model
	AdministratorID uint   `json:"administrator_id" form:"administrator_id"`
	Title           string `json:"title" form:"title"`
	Description     string `json:"description" form:"description"`
	Image           string `json:"image" form:"image"`
	Label           string `json:"label" form:"label"`
	Slug            string `json:"slug" form:"slug"`
}
