package usecase

import (
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/models"
	"go_bedu/repositories"
)

type ArticleUsecase interface {
	GetAllArticles(page, limit int) ([]dtos.ArticleDetailResponse, int, error)
	GetArticleByID(id uint) (dtos.ArticleDetailResponse, error)
	CreateArticle(article *dtos.CreateArticlesRequest) (dtos.ArticleDetailResponse, error)
	UpdateArticle(id uint, article dtos.UpdateArticlesRequest) (dtos.ArticleDetailResponse, error)
	DeleteArticle(id uint) error
}

type articleUsecase struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleRepository(ArticleRepository repositories.ArticleRepository) ArticleUsecase {
	return &articleUsecase{ArticleRepository}
}

// Logic for Get All Article from DB with optional pagination
func (u *articleUsecase) GetAllArticles(page, limit int) ([]dtos.ArticleDetailResponse, int, error) {
	articles, count, err := u.articleRepository.GetAllArticles(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var articleResponses []dtos.ArticleDetailResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, dtos.ArticleDetailResponse{
			ArticleID:   article.ID,
			Title:       article.Title,
			Image:       article.Image,
			Description: article.Description,
			Label:       article.Label,
			Slug:        article.Slug,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	return articleResponses, count, nil
}

// Logic for get Article by ID
func (u *articleUsecase) GetArticleByID(id uint) (dtos.ArticleDetailResponse, error) {
	var articleResponses dtos.ArticleDetailResponse

	article, err := u.articleRepository.GetArticleByID(id)
	if err != nil {
		return articleResponses, err
	}

	articleResponse := dtos.ArticleDetailResponse{
		ArticleID:   article.ID,
		Title:       article.Title,
		Image:       article.Image,
		Description: article.Description,
		Label:       article.Label,
		Slug:        article.Slug,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return articleResponse, nil
}

// Logic for Create Article
func (u *articleUsecase) CreateArticle(article *dtos.CreateArticlesRequest) (dtos.ArticleDetailResponse, error) {
	var articleResponses dtos.ArticleDetailResponse

	slug := helpers.CreateSlug(article.Title)

	CreateArticle := models.Article{
		AdministratorID: article.AdministratorID,
		Title:           article.Title,
		Description:     article.Description,
		Image:           article.Image,
		Label:           article.Label,
		Slug:            slug,
	}

	createdArticle, err := u.articleRepository.CreateArticle(CreateArticle)
	fmt.Println(CreateArticle)
	if err != nil {
		return articleResponses, err
	}

	articleResponse := dtos.ArticleDetailResponse{
		ArticleID:   createdArticle.ID,
		Title:       createdArticle.Title,
		Image:       createdArticle.Image,
		Description: createdArticle.Description,
		Label:       createdArticle.Label,
		Slug:        createdArticle.Slug,
		CreatedAt:   createdArticle.CreatedAt,
		UpdatedAt:   createdArticle.UpdatedAt,
	}

	return articleResponse, nil
}

// Logic for Udpate Article
func (u *articleUsecase) UpdateArticle(id uint, article dtos.UpdateArticlesRequest) (dtos.ArticleDetailResponse, error) {
	var (
		articles        models.Article
		articleResponse dtos.ArticleDetailResponse
	)

	articles, err := u.articleRepository.GetArticleByID(id)
	if err != nil {
		return articleResponse, err
	}

	slug := helpers.CreateSlug(articles.Title)

	articles.Title = article.Title
	articles.Description = article.Description
	articles.Image = article.Image
	articles.Label = article.Label
	articles.Slug = slug

	articles, err = u.articleRepository.UpdateArticle(articles)
	if err != nil {
		return articleResponse, err
	}

	articleResponse.ArticleID = articles.ID
	articleResponse.Title = articles.Title
	articleResponse.Image = articles.Image
	articleResponse.Description = articles.Description
	articleResponse.Label = articles.Label
	articleResponse.Slug = articles.Slug
	articleResponse.CreatedAt = articles.CreatedAt
	articleResponse.UpdatedAt = articles.UpdatedAt

	return articleResponse, nil

}

// Logic For Delete Article from DB
func (u *articleUsecase) DeleteArticle(id uint) error {
	article, err := u.articleRepository.GetArticleByID(id)
	if err != nil {
		return err
	}

	err = u.articleRepository.DeleteArticle(article)
	return err
}
