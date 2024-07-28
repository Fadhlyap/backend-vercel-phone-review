package api

import (
	"backend-vercel-phone-review/config"
	"backend-vercel-phone-review/docs"
	"backend-vercel-phone-review/routes"
	"backend-vercel-phone-review/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	App *gin.Engine
)

func init() {
	App = gin.New()

	environment := utils.Getenv("ENVIRONMENT", "development")

	if environment == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	docs.SwaggerInfo.Title = "Phone Review REST API"
	docs.SwaggerInfo.Description = "This is REST API for Phone Reviews."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = utils.Getenv("HOST", "localhost:8080/api/v1")
	if environment == "development" {
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	err := config.ConnectDataBase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	App = routes.SetupRouter()
}

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	App.ServeHTTP(w, r)
}
