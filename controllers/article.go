package controllers

import (
	"go_bedu/lib/database"
	"go_bedu/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Get All Article from DB
func GetArticlesControllers(c echo.Context) error {
	articles, err := database.GetArticles()

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	if len(articles) == 0 {
		return c.JSON(http.StatusOK, models.ResponseMessage{
			Message: "No article found",
		})
	}

	allarticle := make([]models.ArticlesResponse, len(articles))
	for i, article := range articles {
		allarticle[i] = models.ArticlesResponse{
			Title:     article.Title,
			Content:   article.Content,
			ImageLink: article.ImageLink,
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success get all article",
		Data:    allarticle,
	})
}

// Get Article by ID from DB
func GetArticleByIDController(c echo.Context) error {
	articles, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	article, err := database.GetArticleById(articles)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	ArticleResponse := models.ArticlesResponse{
		Title:     article.Title,
		Content:   article.Content,
		ImageLink: article.ImageLink,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success get article by id",
		Data:    ArticleResponse,
	})
}

// Create Article to DB
func CreateArticleController(c echo.Context) error {
	article := models.Article{}
	c.Bind(&article)

	if err := c.Validate(&article); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	article, err := database.CreateArticle(article)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	ArticleResponse := models.ArticlesResponse{
		Title:     article.Title,
		Content:   article.Content,
		ImageLink: article.ImageLink,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success create article",
		Data:    ArticleResponse,
	})
}

// Update Article by ID to DB
func UpdateArticleByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	article := models.Article{}
	c.Bind(&article)

	if err := c.Validate(&article); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	article, err = database.UpdateArticle(article, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	ArticleResponse := models.ArticlesResponse{
		Title:     article.Title,
		Content:   article.Content,
		ImageLink: article.ImageLink,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success update article",
		Data:    ArticleResponse,
	})
}

// Delete Article by ID from DB
func DeleteArticleController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	// Check Apakah ID ada di DB
	_, err = database.GetArticleById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	_, err = database.DeleteArticle(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseMessage{
		Message: "success delete article",
	})
}
