package database

import (
	"go_bedu/config"
	"go_bedu/models"
)

// Get All Articles from DB
func GetArticles() (articles []models.Article, err error) {
	err = config.DB.Find(&articles).Error

	if err != nil {
		return []models.Article{}, err
	}
	return articles, nil
}

// Create Article to DB
func CreateArticle(article models.Article) (models.Article, error) {
	err := config.DB.Create(&article).Error

	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

// Get Article by ID from DB
func GetArticleById(id int) (article models.Article, err error) {
	err = config.DB.Table("articles").Where("id = ?", id).Find(&article).Error

	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

// Update Article by ID from DB
func UpdateArticle(article models.Article, id int) (models.Article, error) {
	err := config.DB.Table("articles").Where("id = ?", id).Updates(&article).Error

	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

// Delete Article by ID from DB
func DeleteArticle(id int) (interface{}, error) {
	err := config.DB.Where("id = ?", id).Delete(&models.Article{}).Error

	if err != nil {
		return nil, err
	}

	return "Article behasil dihapus", nil
}
