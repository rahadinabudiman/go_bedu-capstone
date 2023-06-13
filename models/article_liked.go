package models

import "gorm.io/gorm"

type ArticleLiked struct {
	gorm.Model
	ArticleID uint `gorm:"not null" json:"article_id" form:"article_id"`
	Article   Article
	UserID    uint `gorm:"not null" json:"user_id" form:"user_id"`
}
