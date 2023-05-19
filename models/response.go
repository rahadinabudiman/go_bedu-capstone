package models

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Response Message Only
type ResponseMessage struct {
	Message string `json:"message"`
}
