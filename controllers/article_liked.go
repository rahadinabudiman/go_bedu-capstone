package controllers

import (
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/models"
	"go_bedu/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleLikedControllers interface {
	GetArticleLikedByUserIdController(c echo.Context) error
	CreateArticleLikedController(c echo.Context) error
	DeleteArticleLikedController(c echo.Context) error
}

type articleLikedControllers struct {
	articleLikedUsecase usecase.ArticleLikedUsecase
	articleUsecase      usecase.ArticleUsecase
}

func NewArticleLikedControllers(articleLikedUsecase usecase.ArticleLikedUsecase, articleUsecase usecase.ArticleUsecase) ArticleLikedControllers {
	return &articleLikedControllers{
		articleLikedUsecase: articleLikedUsecase,
		articleUsecase:      articleUsecase,
	}
}

func (c *articleLikedControllers) GetArticleLikedByUserIdController(ctx echo.Context) error {
	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Please login for access",
				helpers.GetErrorData(err),
			),
		)
	}

	articleLiked, err := c.articleLikedUsecase.GetArticleLikedByUserId(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get Article Liked",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Get Article Liked",
			articleLiked,
		),
	)
}

func (c *articleLikedControllers) CreateArticleLikedController(ctx echo.Context) error {
	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Please login for access",
				helpers.GetErrorData(err),
			),
		)
	}

	idArticle, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot get article",
				helpers.GetErrorData(err),
			),
		)
	}

	dataArticle, err := c.articleUsecase.GetArticleByID(uint(idArticle))
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed to get article",
				helpers.GetErrorData(err),
			),
		)
	}

	if dataArticle.ArticleID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Article not found",
				nil,
			),
		)
	}

	article := models.ArticleLiked{
		UserID:    uint(id),
		ArticleID: dataArticle.ArticleID,
	}

	articleLiked, err := c.articleLikedUsecase.CreateArticleLiked(uint(id), article)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot bookmark article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Article has been saved",
			articleLiked,
		),
	)

	return nil
}

func (c *articleLikedControllers) DeleteArticleLikedController(ctx echo.Context) error {
	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Please login for access",
				helpers.GetErrorData(err),
			),
		)
	}

	idArticle, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot get article",
				helpers.GetErrorData(err),
			),
		)
	}

	DataArticle, err := c.articleUsecase.GetArticleByID(uint(idArticle))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot get article",
				helpers.GetErrorData(err),
			),
		)
	}

	UserIDData := DataArticle.ArticleID

	Article := models.ArticleLiked{
		UserID:    uint(id),
		ArticleID: uint(UserIDData),
	}

	err = c.articleLikedUsecase.DeleteArticleLiked(uint(id), Article)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot Delete Bookmark Article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponseMessage(
			http.StatusOK,
			"Article has been deleted",
		),
	)
}
