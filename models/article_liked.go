package models

import "gorm.io/gorm"

type ArticleLiked struct {
	gorm.Model
	ArticleID uint    `json:"article_id" form:"article_id"`
	Article   Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID    uint    `json:"user_id" form:"user_id"`
}
