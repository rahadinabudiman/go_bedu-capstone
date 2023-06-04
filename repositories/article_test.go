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
	})
}
