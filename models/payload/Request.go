package payload

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterAdminRequest struct {
	Nama     string `json:"nama" form:"nama" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"gte=6"`
	Role     string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null"`
}

type UpdateAdminRequest struct {
	Nama     string `json:"nama" form:"nama" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"gte=6"`
	Role     string `json:"role" form:"role" gorm:"type:enum('Admin', 'Super Admin');default:'Admin'; not-null"`
}

type DeleteAdminRequest struct {
	Password string `json:"password" form:"password" validate:"gte=6"`
}

type CreateArticlesRequest struct {
	IDAdmin   uint   `json:"id_admin" form:"id_admin"`
	Title     string `json:"title" form:"title" validate:"required"`
	Content   string `json:"content" form:"content" validate:"required"`
	ImageLink string `json:"image_link" form:"image_link" validate:"required"`
}

type UpdateArticlesRequest struct {
	Title     string `json:"title" form:"title" validate:"required"`
	Content   string `json:"content" form:"content" validate:"required"`
	ImageLink string `json:"image_link" form:"image_link" validate:"required"`
}
