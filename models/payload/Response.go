package payload

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Response Message Only
type ResponseMessage struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Email string `json:"email" form:"email"`
	Token string `json:"token" form:"token"`
}

type RegisterAdminResponse struct {
	Nama     string `json:"nama" form:"nama"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

type UpdateAdminResponse struct {
	Nama     string `json:"nama" form:"nama"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}

type CreateArticlesResponse struct {
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	ImageLink string `json:"image_link" form:"image_link"`
}

type UpdateArticleResponse struct {
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	ImageLink string `json:"image_link" form:"image_link"`
}

type AdminProfileResponse struct {
	ID    uint   `json:"id" form:"id"`
	Nama  string `json:"nama" form:"nama"`
	Email string `json:"email" form:"email"`
	Role  string `json:"role" form:"role"`
}

type ArticleDetailResponse struct {
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	ImageLink string `json:"image_link" form:"image_link"`
}
