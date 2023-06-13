package repositories

import (
	"go_bedu/models"

	"gorm.io/gorm"
)

type ArticleLikedRepository interface {
	GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error)
	GetLikeByUserIdAndArticleId(userId uint, articleId uint) (models.ArticleLiked, error)
	CreateArticleLiked(articleLiked models.ArticleLiked) (models.ArticleLiked, error)
	DeleteArticleLiked(userId uint, articleId uint) (articleLiked models.ArticleLiked, err error)
}

type articleLikedRepository struct {
	db *gorm.DB
}

func NewArticleLikedRepository(db *gorm.DB) *articleLikedRepository {
	return &articleLikedRepository{db}
}

func (r *articleLikedRepository) GetLikeByUserIdAndArticleId(userId uint, articleId uint) (models.ArticleLiked, error) {
	var articleLiked models.ArticleLiked

	err := r.db.Where("user_id = ? AND article_id = ?", userId, articleId).First(&articleLiked).Error

	return articleLiked, err
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

func (r *articleLikedRepository) DeleteArticleLiked(userId uint, articleId uint) (articleLiked models.ArticleLiked, err error) {
	err = r.db.Where("user_id = ? AND article_id = ?", userId, articleId).Delete(&articleLiked).Error

	return articleLiked, err
}
