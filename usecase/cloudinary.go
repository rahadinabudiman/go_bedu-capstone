package usecase

import (
	"go_bedu/helpers"
	"go_bedu/models"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type ClourdinaryUsecase interface {
	FileUpload(file models.File) (string, error)
	RemoteUpload(url models.Url) (string, error)
}

type media struct{}

func NewMediaUpload() ClourdinaryUsecase {
	return &media{}
}

// FileUpload godoc
// @Summary      Upload file
// @Description  Upload file to cloudinary
// @Tags         Cloudinary
// @Accept       json
// @Produce      json
// @Param        file formData file false "Photo file"
// @Success      200 {object} dtos.StatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/cloudinary/file-upload [post]
func (*media) FileUpload(file models.File) (string, error) {
	//validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := helpers.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

// RemoteUpload godoc
// @Summary      Upload file
// @Description  Upload file to cloudinary
// @Tags         Cloudinary
// @Accept       json
// @Produce      json
// @Param        request body models.Url true "Payload Body [RAW]"
// @Success      200 {object} dtos.StatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/cloudinary/url-upload [post]
func (*media) RemoteUpload(url models.Url) (string, error) {
	//validate
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, errUrl := helpers.ImageUploadHelper(url.Url)
	if errUrl != nil {
		return "", err
	}
	return uploadUrl, nil
}
