package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model `swaggerignore:"true"`
	PhoneID    uint      `json:"phone_id"`
	UserID     uint      `json:"user_id"`
	Rating     int       `json:"rating"`
	Content    string    `json:"content"`
	Comments   []Comment `json:"comments" gorm:"foreignKey:ReviewID" swaggerignore:"true"`
}
