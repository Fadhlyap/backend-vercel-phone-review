package models

import "gorm.io/gorm"

type Phone struct {
	gorm.Model `swaggerignore:"true"`
	Name       string    `json:"name"`
	Brand      string    `json:"brand"`
	Features   []Feature `json:"features" gorm:"foreignKey:PhoneID"`
	Reviews    []Review  `json:"reviews" gorm:"foreignKey:PhoneID"`
}
