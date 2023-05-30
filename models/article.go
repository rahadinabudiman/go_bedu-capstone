package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	AdministratorID uint   `json:"administrator_id" form:"administrator_id"`
	Thumbnail       string `json:"thumbnail" form:"thumbnail"`
	Title           string `json:"title" form:"title"`
	Abstract        string `json:"abstract" form:"abstract"`
	Description     string `json:"description" form:"description"`
	Image           string `json:"image" form:"image"`
	Label           string `json:"label" form:"label"`
	Slug            string `json:"slug" form:"slug"`
}
