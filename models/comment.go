package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model `swaggerignore:"true"`
	ReviewID   uint   `json:"review_id"`
	UserID     uint   `json:"user_id"`
	Content    string `json:"content"`
}
