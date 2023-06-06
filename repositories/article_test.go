package repositories

import (
	"go_bedu/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestArticleRepository(t *testing.T) {
	// Setup
	dsn := "r4ha:kmoonkinan@tcp(localhost:3306)/go_bedu?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	articleRepository := NewArticleRepository(db)
	t.Run("Test Create Article", func(t *testing.T) {
		// Get Last Admin ID
		lastAdmin := models.Administrator{}
		result := db.Last(&lastAdmin)
		assert.NoError(t, result.Error)

		article := models.Article{
			AdministratorID: lastAdmin.ID,
			Thumbnail:       "thumbnail.jpg",
			Title:           "Test Article",
			Abstract:        "This is a test article",
			Description:     "Lorem ipsum dolor sit amet",
			Image:           "image.jpg",
			Label:           "Test",
			Slug:            "test-article",
		}

		// Create Article
		CreatedArticle, err := articleRepository.CreateArticle(article)
		assert.NoError(t, err)
		assert.NotZero(t, CreatedArticle.ID)

		// Check if the article exist already
		dbArticle, err := articleRepository.GetArticleByID(uint(CreatedArticle.ID))
		assert.NoError(t, err)
		assert.Equal(t, CreatedArticle.ID, dbArticle.ID)
	})

	t.Run("Update Article", func(t *testing.T) {
		// Get Last Article ID
		lastArticle := models.Article{}
		result := db.Last(&lastArticle)
		assert.NoError(t, result.Error)

		lastArticle.Thumbnail = "thumbnailupdate.jpg"
		lastArticle.Title = "Test Update Article"
		lastArticle.Abstract = "This is an updated article"
		lastArticle.Description = "Updated Lorem ipsum dolor sit amet"
		lastArticle.Image = "imageupdate.jpg"
		lastArticle.Label = "Test Update"
		lastArticle.Slug = "test-update"

		updated, err := articleRepository.UpdateArticle(lastArticle)
		assert.NoError(t, err)

		// Verify Updated Article
		assert.Equal(t, lastArticle.Slug, updated.Slug)
	})

	t.Run("Test Delete Article", func(t *testing.T) {
		// Get Last Article ID
		lastArticle := models.Article{}
		result := db.Last(&lastArticle)
		assert.NoError(t, result.Error)

		// Get Article By ID
		article, err := articleRepository.GetArticleByID(uint(lastArticle.ID))
		assert.NoError(t, err)

		// Delete Article
		err = articleRepository.DeleteArticle(article)
		assert.NoError(t, err)

		// Verify Article was Deleted
		DeletedArticle, err := articleRepository.GetArticleByID(uint(article.ID))

		// Check if the deleted article is empty
		assert.Error(t, err)
		assert.Equal(t, models.Article{}, DeletedArticle)
	})
}
