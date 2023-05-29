package dtos

import "go_bedu/helpers"

type LoginStatusOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Email      string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
	Password   string `json:"password" form:"password" validate:"required" example:"rahadinabudimansundara"`
}

type AdminStatusOKResponse struct {
	StatusCode int                 `json:"status_code" example:"200"`
	Message    string              `json:"message" example:"Successfully get user credentials"`
	Data       AdminDetailResponse `json:"data"`
}

type AdminCreeatedResponse struct {
	StatusCode int                 `json:"status_code" example:"201"`
	Message    string              `json:"message" example:"Successfully registered"`
	Data       AdminDetailResponse `json:"data"`
}

type GetAllAdminsResponse struct {
	StatusCode int                 `json:"status_code" example:"201"`
	Message    string              `json:"message" example:"Successfully registered"`
	Data       AdminDetailResponse `json:"data"`
}

type GetAllArticleStatusOKResponse struct {
	StatusCode int                   `json:"status_code" example:"200"`
	Message    string                `json:"message" example:"Successfully get article"`
	Data       ArticleDetailResponse `json:"data"`
	Meta       helpers.Meta          `json:"meta"`
}

type ArticleStatusOKResponse struct {
	StatusCode int                   `json:"status_code" example:"200"`
	Message    string                `json:"message" example:"Successfully get article"`
	Data       ArticleDetailResponse `json:"data"`
}

type ArticleCreeatedResponse struct {
	StatusCode int                   `json:"status_code" example:"201"`
	Message    string                `json:"message" example:"Successfully created article"`
	Data       ArticleDetailResponse `json:"data"`
}

type StatusOKDeletedResponse struct {
	StatusCode int         `json:"status_code" example:"200"`
	Message    string      `json:"message" example:"Successfully deleted"`
	Errors     interface{} `json:"errors"`
}

type BadRequestResponse struct {
	StatusCode int         `json:"status_code" example:"400"`
	Message    string      `json:"message" example:"Bad Request"`
	Errors     interface{} `json:"errors"`
}

type UnauthorizedResponse struct {
	StatusCode int         `json:"status_code" example:"401"`
	Message    string      `json:"message" example:"Unauthorized"`
	Errors     interface{} `json:"errors"`
}

type ForbiddenResponse struct {
	StatusCode int         `json:"status_code" example:"403"`
	Message    string      `json:"message" example:"Forbidden"`
	Errors     interface{} `json:"errors"`
}

type NotFoundResponse struct {
	StatusCode int         `json:"status_code" example:"404"`
	Message    string      `json:"message" example:"Not Found"`
	Errors     interface{} `json:"errors"`
}

type InternalServerErrorResponse struct {
	StatusCode int         `json:"status_code" example:"500"`
	Message    string      `json:"message" example:"Internal Server Error"`
	Errors     interface{} `json:"errors"`
}
