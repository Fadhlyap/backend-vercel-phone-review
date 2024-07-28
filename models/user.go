package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `swaggerignore:"true"`
	Username   string    `json:"username" gorm:"unique"`
	Password   string    `json:"-"`
	Profile    Profile   `json:"profile" gorm:"foreignkey:UserID"`
	Reviews    []Review  `json:"reviews"`
	Comment    Comment `json:"comments"`
}
