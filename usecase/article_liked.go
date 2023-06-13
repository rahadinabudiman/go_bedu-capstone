package usecase

import (
	"go_bedu/models"
	"go_bedu/repositories"

	"github.com/labstack/echo/v4"
)

type ArticleLikedUsecase interface {
	GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error)
	GetArticleLikeByUserIdAndArticleId(userId uint, articleId uint) (models.ArticleLiked, error)
	CreateArticleLiked(id uint, articleId uint) (models.ArticleLiked, error)
	DeleteArticleLiked(id uint, articleId uint) (articleLiked models.ArticleLiked, err error)
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

func (u *articleLikedUsecase) GetArticleLikeByUserIdAndArticleId(userId uint, articleId uint) (models.ArticleLiked, error) {
	var articlelike models.ArticleLiked

	articleLiked, err := u.articleLikedRepo.GetLikeByUserIdAndArticleId(userId, articleId)
	if err != nil {
		return articlelike, err
	}

	return articleLiked, nil
}

func (u *articleLikedUsecase) CreateArticleLiked(id uint, articleId uint) (models.ArticleLiked, error) {
	articleLiked := models.ArticleLiked{
		ArticleID: articleId,
		UserID:    id,
	}

	CountArticle, _ := u.articleLikedRepo.GetLikeByUserIdAndArticleId(uint(id), uint(articleId))
	if CountArticle.ID != 0 {
		articleLiked, _ = u.articleLikedRepo.DeleteArticleLiked(id, articleId)
	} else {
		articleLiked, _ = u.articleLikedRepo.CreateArticleLiked(articleLiked)
	}

	return articleLiked, nil
}

func (u *articleLikedUsecase) DeleteArticleLiked(id uint, articleId uint) (articleLiked models.ArticleLiked, err error) {
	_, err = u.userRepository.ReadToken(id)
	if err != nil {
		echo.NewHTTPError(400, "Failed to get User")
	}

	articleLiked, err = u.articleLikedRepo.DeleteArticleLiked(id, articleId)
	if err != nil {
		return articleLiked, err
	}

	return articleLiked, nil
}
