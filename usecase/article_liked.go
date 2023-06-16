package usecase

import (
	"go_bedu/models"
	"go_bedu/repositories"
)

type ArticleLikedUsecase interface {
	GetArticleLikedByUserId(userId uint) ([]models.ArticleLiked, error)
	GetArticleLikeByUserIdAndArticleId(userId uint, articleId uint) (models.ArticleLiked, error)
	CreateArticleLiked(id uint, articleId uint) (models.ArticleLiked, error)
}

type articleLikedUsecase struct {
	articleLikedRepo repositories.ArticleLikedRepository
	userRepository   repositories.UserRepository
}

func NewArticleLikedUsecase(articleLikedRepo repositories.ArticleLikedRepository, userRepository repositories.UserRepository) *articleLikedUsecase {
	return &articleLikedUsecase{articleLikedRepo, userRepository}
}

// GetBookmarkArticle godoc
// @Summary      Get Bookmark by User ID
// @Description  Get Bookmark by User ID
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param id path integer true "ID user"
// @Success      200 {object} dtos.LikedStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/liked/{id} [get]
// @Security     BearerAuth
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

// CreateBookmarkArticle godoc
// @Summary      Create Bookmark by Article ID
// @Description  Create Bookmark by Article ID
// @Tags         Article
// @Accept       json
// @Produce      json
// @Param id path integer true "ID article"
// @Success      200 {object} dtos.LikedStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /article/like/{id} [get]
// @Security     BearerAuth
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
