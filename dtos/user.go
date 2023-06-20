package dtos

import "time"

type LogoutUserResponse struct {
	Message string `json:"message" form:"message" example:"Logout Success"`
}

type UserDetailResponse struct {
	ID        uint      `json:"id" from:"id"`
	Username  string    `json:"username" form:"username" validate:"required" example:"r4ha"`
	Nama      string    `json:"nama" from:"nama"`
	Email     string    `json:"email" from:"email"`
	Role      string    `json:"role" from:"role"`
	CreatedAt time.Time `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
	// Article   []models.Article `json:"article" from:"article"`
}

type RegisterUserRequest struct {
	Nama            string `json:"nama" form:"nama" validate:"required" example:"Rahadina Budiman Sundara"`
	Username        string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email           string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
	Password        string `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
	PasswordConfirm string `json:"passwordconfirm" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
	Verified        bool   `gorm:"type:enum('False', 'True');default:'False'; not-null" example:"False"`
	Role            string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null" example:"Admin"`
}

type UpdateUserRequest struct {
	Nama     string `json:"nama" form:"nama" validate:"required" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email    string `json:"email" form:"email" validate:"required,email" example:"me@r4ha.com"`
	Role     string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null" example:"Admin"`
}

type DeleteUserRequest struct {
	Password string `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
}

type RegisterUserResponse struct {
	Nama     string `json:"nama" form:"nama" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email    string `json:"email" form:"email" example:"me@r4ha.com"`
	Password string `json:"password" form:"password" example:"rahadinabudimansundara"`
	Role     string `json:"role" form:"role" example:"Admin"`
}

type UpdateUserResponse struct {
	Nama     string `json:"nama" form:"nama" example:"Rahadina Budiman Sundara"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Email    string `json:"email" form:"email" example:"me@r4ha.com"`
	Role     string `json:"role" form:"role" example:"Admin"`
}

type UserProfileResponse struct {
	ID       uint   `json:"id" form:"id" example:"1"`
	Username string `json:"username" form:"username" validate:"required" example:"r4ha"`
	Nama     string `json:"nama" form:"nama" example:"Rahadina Budiman Sundara"`
	Email    string `json:"email" form:"email" example:"me@r4ha.com"`
	Role     string `json:"role" form:"role" example:"Admin"`
}

type ChangePasswordUserRequest struct {
	OldPassword     string `json:"old_password" form:"old_password" validate:"gte=6" example:"rahadinabudimansundara"`
	Password        string `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
	PasswordConfirm string `json:"passwordconfirm" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
}
