package dtos

type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Password string `json:"password" form:"password" validate:"required" example:"rahadinabudimansundara"`
}

type LoginResponse struct {
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Token    string `json:"token" form:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type VerifyEmailResponse struct {
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Message  string `json:"message" form:"message" example:"Email has been verified"`
}

type VerifyEmailRequest struct {
	VerificationCode string `json:"verification_code" form:"verification_code" validate:"required" example:"1234567890"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
}

type ForgotPasswordResponse struct {
	Email   string `json:"email" form:"email" example:"me@r4ha.com"`
	Message string `json:"message" form:"message" example:"Email has been sent"`
}

type ChangePasswordRequest struct {
	Password        string `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
	PasswordConfirm string `json:"passwordconfirm" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
}

type ChangePasswordByOTPResponse struct {
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Message  string `json:"message" form:"message" example:"Password has been reset successfully"`
}
