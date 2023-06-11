package repositories

import (
	"go_bedu/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAllArticles(page, limit int) ([]models.Article, int, error)
	GetArticleByID(id uint) (models.Article, error)
	GetArticleByImage(image string) (int64, error)
	GetArticleByThumbnail(thumbnail string) (int64, error)
	CreateArticle(article models.Article) (models.Article, error)
	UpdateArticle(article models.Article) (models.Article, error)
	DeleteArticle(article models.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

// Get All Articles from DB with optional pagination
func (r *articleRepository) GetAllArticles(page, limit int) ([]models.Article, int, error) {
	var (
		articles []models.Article
		count    int64
	)

	err := r.db.Find(&articles).Count(&count).Error
	if err != nil {
		return articles, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&articles).Error

	return articles, int(count), err
}

// Get Article By ID from DB
func (r *articleRepository) GetArticleByID(id uint) (models.Article, error) {
	var article models.Article

	err := r.db.Where("id = ?", id).First(&article).Error
	return article, err
}

// Get Article by Image to validate and check the images is changes or not
func (r *articleRepository) GetArticleByImage(image string) (int64, error) {
	var (
		totalImage int64
		Article    models.Article
	)
	err := r.db.Model(&Article).Where("image = ?", image).Count(&totalImage).Error
	if err != nil {
		return 0, err
	}

	return totalImage, nil
}

// Get Article by Thumbnails
func (r *articleRepository) GetArticleByThumbnail(thumbnail string) (int64, error) {
	var (
		totalImage int64
		Article    models.Article
	)
	err := r.db.Model(&Article).Where("image = ?", thumbnail).Count(&totalImage).Error
	if err != nil {
		return 0, err
	}

	return totalImage, nil
}

// Create Article and save to DB
func (r *articleRepository) CreateArticle(article models.Article) (models.Article, error) {
	err := r.db.Create(&article).Error

	return article, err
}

// Update Article and save to DB
func (r *articleRepository) UpdateArticle(article models.Article) (models.Article, error) {
	err := r.db.Table("articles").Save(&article).Error

	return article, err
}

// Delete Article from DB
func (r *articleRepository) DeleteArticle(article models.Article) error {
	err := r.db.Delete(&article).Error

	return err
}
