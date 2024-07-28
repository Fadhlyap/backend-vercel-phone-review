package models

import "gorm.io/gorm"

type Feature struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name"`
	Details    string `json:"details"`
	PhoneID    uint   `json:"phone_id" swaggerignore:"true"`
}
