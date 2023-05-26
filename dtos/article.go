package dtos

import "time"

type CreateArticlesRequest struct {
	AdministratorID uint   `json:"administrator_id" form:"administrator_id"`
	Title           string `json:"title" form:"title"`
	Description     string `json:"description" form:"description"`
	Image           string `json:"image" form:"image"`
	Label           string `json:"label" form:"label"`
}

type UpdateArticlesRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	Label       string `json:"label" form:"label"`
}

type CreateArticlesResponse struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	Label       string `json:"label" form:"label"`
	Slug        string `json:"slug" form:"slug"`
}

type UpdateArticleResponse struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
	Label       string `json:"label" form:"label"`
}

type ArticleDetailResponse struct {
	ArticleID   uint      `json:"article_id"`
	Title       string    `json:"title" `
	Image       string    `json:"image" `
	Description string    `json:"description" `
	Label       string    `json:"label" `
	Slug        string    `json:"slug" `
	CreatedAt   time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
