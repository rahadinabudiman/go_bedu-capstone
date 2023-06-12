package controllers

import (
	"go_bedu/helpers"
	"go_bedu/models"
	"go_bedu/usecase"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

type CloudinaryController struct {
	cloudinaryUsecase usecase.ClourdinaryUsecase
}

func NewCloudinaryController(cloudinaryUsecase usecase.ClourdinaryUsecase) CloudinaryController {
	return CloudinaryController{cloudinaryUsecase}
}

func (c *CloudinaryController) FileUpload(ctx echo.Context) error {
	var cloudinary models.Url

	formHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Error uploading photo",
				helpers.GetErrorData(err),
			),
		)
	}

	//get file from header
	formFile, err := formHeader.Open()
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Error uploading photo",
				helpers.GetErrorData(err),
			),
		)
	}

	var re = regexp.MustCompile(`.png|.jpeg|.jpg`)

	if !re.MatchString(formHeader.Filename) {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"The provided file format is not allowed. Please upload a JPEG or PNG image",
				helpers.GetErrorData(err),
			),
		)
	}

	uploadUrl, err := usecase.NewMediaUpload().FileUpload(models.File{File: formFile})

	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Error uploading photo",
				helpers.GetErrorData(err),
			),
		)
	}
	cloudinary.Url = uploadUrl

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully upload file to cloudinary",
			cloudinary,
		),
	)
}

func (c *CloudinaryController) UrlUpload(ctx echo.Context) error {
	var cloudinaryDTO models.Url
	if err := ctx.Bind(&cloudinaryDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding cloudinary",
				helpers.GetErrorData(err),
			),
		)
	}

	cloudinary, err := c.cloudinaryUsecase.RemoteUpload(cloudinaryDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed upload file to cloudinary",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully upload file to cloudinary",
			cloudinary,
		),
	)
}
