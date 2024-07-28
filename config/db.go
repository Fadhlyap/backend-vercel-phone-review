package config

import (
	"backend-vercel-phone-review/models"
	"backend-vercel-phone-review/utils"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() error {
	dbProvider := utils.Getenv("DB_PROVIDER", "mysql")

	var db *gorm.DB
	var err error

	if dbProvider == "postgres" {
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_NAME")
		// production
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", host, username, password, database, port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error connecting to PostgreSQL: %v", err)
			return err
		}
	} else {
		username := utils.Getenv("DB_USERNAME", "root")
		password := utils.Getenv("DB_PASSWORD", "root")
		host := utils.Getenv("DB_HOST", "127.0.0.1")
		port := utils.Getenv("DB_PORT", "3306")
		database := utils.Getenv("DB_NAME", "db_name")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error connecting to MySQL: %v", err)
			return err
		}
	}

	// Set the global DB variable
	DB = db

	// Auto migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Comment{}, &models.Phone{}, &models.Feature{}, &models.Review{})
	if err != nil {
		log.Printf("Error during migration: %v", err)
		return err
	}

	return nil
}
