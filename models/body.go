package models

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PhoneRequest struct {
	Brand string `json:"brand" binding:"required"`
	Name  string `json:"name" binding:"required"`
}
