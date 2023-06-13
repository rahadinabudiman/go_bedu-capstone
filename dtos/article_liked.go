package dtos

type ArticleLiked struct {
	ArticleID uint `json:"article_id" form:"article_id"`
	UserID    uint `json:"user_id" form:"user_id"`
}
