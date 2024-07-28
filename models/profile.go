package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model `swaggerignore:"true"`
	UserID     uint   `json:"user_id"`
	FullName   string `json:"full_name"`
	Bio        string `json:"bio"`
}
