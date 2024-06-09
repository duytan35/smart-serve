package controllers

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationRequest struct {
	Page int `json:"page" required:"true" example:"1"`
	Size int `json:"size" required:"true" example:"10"`
}

type PaginationMetaResponse struct {
	Page      int `json:"page" required:"true" example:"1"`
	Size      int `json:"size" required:"true" example:"10"`
	TotalPage int `json:"totalPage" required:"true" example:"2"`
}

type PaginationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Meta    PaginationMetaResponse
	Data    interface{} `json:"data,omitempty"`
}
