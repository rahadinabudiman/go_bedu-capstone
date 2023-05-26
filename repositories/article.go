package repositories

import (
	"go_bedu/models"

	"github.com/jinzhu/gorm"
)

type ArticleRepository interface {
	GetAllArticles(page, limit int) ([]models.Article, int, error)
	GetArticleByID(id uint) (models.Article, error)
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

// Create Article and save to DB
func (r *articleRepository) CreateArticle(article models.Article) (models.Article, error) {
	err := r.db.Create(&article).Error

	return article, err
}

// Update Article and save to DB
func (r *articleRepository) UpdateArticle(article models.Article) (models.Article, error) {
	err := r.db.Save(&article).Error

	return article, err
}

// Delete Article from DB
func (r *articleRepository) DeleteArticle(article models.Article) error {
	err := r.db.Delete(&article).Error

	return err
}
