package usecase

import (
	"go_bedu/models"
	"go_bedu/repositories"

	"github.com/labstack/echo/v4"
)

type ArticleLikedUsecase interface {
	GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error)
	CreateArticleLiked(id uint, articleLiked models.ArticleLiked) (models.ArticleLiked, error)
	DeleteArticleLiked(id uint, articleLiked models.ArticleLiked) error
}

type articleLikedUsecase struct {
	articleLikedRepo repositories.ArticleLikedRepository
	userRepository   repositories.UserRepository
}

func NewArticleLikedUsecase(articleLikedRepo repositories.ArticleLikedRepository, userRepository repositories.UserRepository) *articleLikedUsecase {
	return &articleLikedUsecase{articleLikedRepo, userRepository}
}

func (u *articleLikedUsecase) GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error) {
	user, err := u.userRepository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	articleLiked, err := u.articleLikedRepo.GetArticleLikedByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	return articleLiked, nil
}

func (u *articleLikedUsecase) CreateArticleLiked(id uint, articleLiked models.ArticleLiked) (models.ArticleLiked, error) {
	_, err := u.userRepository.ReadToken(id)
	if err != nil {
		echo.NewHTTPError(400, "Failed to get User")
	}

	articleLiked, err = u.articleLikedRepo.CreateArticleLiked(articleLiked)
	if err != nil {
		return articleLiked, err
	}

	return articleLiked, nil
}

func (u *articleLikedUsecase) DeleteArticleLiked(id uint, articleLiked models.ArticleLiked) error {
	_, err := u.userRepository.ReadToken(id)
	if err != nil {
		echo.NewHTTPError(400, "Failed to get User")
	}

	err = u.articleLikedRepo.DeleteArticleLiked(articleLiked)
	if err != nil {
		return err
	}

	return nil
}
