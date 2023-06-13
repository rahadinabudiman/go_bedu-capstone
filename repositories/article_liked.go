package repositories

import (
	"go_bedu/models"

	"gorm.io/gorm"
)

type ArticleLikedRepository interface {
	GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error)
	CreateArticleLiked(articleLiked models.ArticleLiked) (models.ArticleLiked, error)
	DeleteArticleLiked(articleLiked models.ArticleLiked) error
}

type articleLikedRepository struct {
	db *gorm.DB
}

func NewArticleLikedRepository(db *gorm.DB) *articleLikedRepository {
	return &articleLikedRepository{db}
}

func (r *articleLikedRepository) GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error) {
	var articleLiked []models.ArticleLiked

	err := r.db.Preload("Article").Where("user_id = ?", userId).Find(&articleLiked).Error

	return articleLiked, err
}

func (r *articleLikedRepository) CreateArticleLiked(articleLiked models.ArticleLiked) (models.ArticleLiked, error) {
	err := r.db.Create(&articleLiked).Error

	return articleLiked, err
}

func (r *articleLikedRepository) DeleteArticleLiked(articleLiked models.ArticleLiked) error {
	err := r.db.Delete(&articleLiked).Error

	return err
}
