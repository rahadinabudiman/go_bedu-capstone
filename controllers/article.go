package controllers

import (
	"go_bedu/dtos"
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/usecase"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ArticleController interface {
	GetAllArticles(c echo.Context) error
	GetArticleById(c echo.Context) error
	CreateArticle(c echo.Context) error
	UpdateArticle(c echo.Context) error
	DeleteArticle(c echo.Context) error
}

type articleController struct {
	articleUsecase usecase.ArticleUsecase
}

func NewArticleController(articleUsecase usecase.ArticleUsecase) ArticleController {
	return &articleController{articleUsecase}
}

// Controller for Get All Article from DB with optional pagination
func (c *articleController) GetAllArticles(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	articles, count, err := c.articleUsecase.GetAllArticles(page, limit)
	if err != nil {

		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching articles",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all article",
			articles,
			page,
			limit,
			count,
		),
	)
}

// Controller for get Article by ID from parameter
func (c *articleController) GetArticleById(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article",
				helpers.GetErrorData(err),
			),
		)
	}

	article, err := c.articleUsecase.GetArticleByID(uint(id))
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

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully get article",
			article,
		),
	)
}

func (c *articleController) CreateArticle(ctx echo.Context) error {
	var articleInput dtos.CreateArticlesRequest
	// Get Admin id from JWT Cookie
	id, _ := m.IsAdmin(ctx)
	articleInput.AdministratorID = uint(id)

	if err := ctx.Bind(&articleInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to bind article",
				helpers.GetErrorData(err),
			),
		)
	}

	// Upload File and validate file extension (jpg, png, and jpeg).
	thumbnail, err := ctx.FormFile("thumbnail")
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to upload thumbnail",
				helpers.GetErrorData(err),
			),
		)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to upload image",
				helpers.GetErrorData(err),
			),
		)
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to open file",
				helpers.GetErrorData(err),
			),
		)
	}

	re := regexp.MustCompile(`.png|.jpeg|.jpg`)

	if !re.MatchString(file.Filename) {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewResponseMessage(
				http.StatusBadRequest,
				"The provided file format is not allowed. Please upload a JPEG or PNG image",
			),
		)
	}

	defer src.Close()

	imageExtension := helpers.GetFileExtension(file.Filename)
	thumbnailExtension := helpers.GetFileExtension(thumbnail.Filename)

	// Generate new filenames
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	newImageFilename := timestamp + imageExtension
	newThumbnailFilename := timestamp + thumbnailExtension

	// Destination paths
	imageDst := "public/images/" + newImageFilename
	thumbnailDst := "public/images/" + newThumbnailFilename

	// Create and save image file
	imageDstFile, err := os.Create(imageDst)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to upload image",
				helpers.GetErrorData(err),
			),
		)
	}
	defer imageDstFile.Close()

	// Copy image file
	if _, err = io.Copy(imageDstFile, src); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to copy image",
				helpers.GetErrorData(err),
			),
		)
	}

	// Create and save thumbnail file
	thumbnailSrc, err := thumbnail.Open()
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to open thumbnail file",
				helpers.GetErrorData(err),
			),
		)
	}
	defer thumbnailSrc.Close()

	thumbnailDstFile, err := os.Create(thumbnailDst)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to upload thumbnail",
				helpers.GetErrorData(err),
			),
		)
	}
	defer thumbnailDstFile.Close()

	// Copy thumbnail file
	if _, err = io.Copy(thumbnailDstFile, thumbnailSrc); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to copy thumbnail",
				helpers.GetErrorData(err),
			),
		)
	}

	articleInput.Image = newImageFilename
	articleInput.Thumbnail = newThumbnailFilename

	article, err := c.articleUsecase.CreateArticle(&articleInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to create article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully create article",
			article,
		),
	)
}

// Controller for update Article by ID from Param
func (c *articleController) UpdateArticle(ctx echo.Context) error {
	var articleInput dtos.UpdateArticlesRequest
	// Get Admin id from JWT Cookie
	id, _ := m.IsAdmin(ctx)
	articleInput.AdministratorID = uint(id)

	if err := ctx.Bind(&articleInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to binding article",
				helpers.GetErrorData(err),
			),
		)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article",
				helpers.GetErrorData(err),
			),
		)
	}

	article, err := c.articleUsecase.GetArticleByID(uint(id))
	if article.ArticleID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article",
				helpers.GetErrorData(err),
			),
		)
	}

	articleRespon, err := c.articleUsecase.UpdateArticle(uint(id), articleInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to update article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Article updated successfully",
			articleRespon,
		),
	)
}

// Controller for delete article by id from params
func (c *articleController) DeleteArticle(ctx echo.Context) error {
	_, err := m.IsAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for Admin Only",
				helpers.GetErrorData(err),
			),
		)
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get article ID",
				helpers.GetErrorData(err),
			),
		)
	}

	err = c.articleUsecase.DeleteArticle(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete article",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponseMessage(
			http.StatusOK,
			"Article deleted successfully",
		),
	)
}
