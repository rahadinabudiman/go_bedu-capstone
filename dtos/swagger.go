package dtos

import "go_bedu/helpers"

type LogoutAdminOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" form:"message" example:"Logout Success"`
}

type VerifyEmailOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Username   string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Message    string `json:"message" form:"message" example:"Email has been verified"`
}

type ForgotPasswordOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Email      string `json:"email" form:"email" example:"me@r4ha.com"`
	Message    string `json:"message" form:"message" example:"OTP has been sent to your email"`
}
type ChangePasswordOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Email      string `json:"email" form:"email" example:"me@r4ha.com"`
	Message    string `json:"message" form:"message" example:"Password has been reset successfully"`
}

type ChangePasswordAdminOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" form:"message" example:"Password has been reset successfully"`
}

type LoginStatusOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Username   string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Token      string `json:"token" form:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
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
