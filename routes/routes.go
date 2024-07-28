package routes

import (
	"backend-vercel-phone-review/controllers"
	"backend-vercel-phone-review/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", controllers.Register)
			authRoutes.POST("/login", controllers.Login)
			authRoutes.GET("/me", middleware.JWTAuthMiddleware(), controllers.GetMe)
			authRoutes.PUT("/change-password/:id", controllers.ChangePassword)
		}

		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.JWTAuthMiddleware())
		{
			userRoutes.GET("/:id", controllers.GetUser)
			userRoutes.PUT("/:id/profile", controllers.UpdateProfile)
		}

		phoneRoutes := api.Group("/phones")
		{
			phoneRoutes.GET("/", controllers.GetPhones)
			phoneRoutes.POST("/", middleware.JWTAuthMiddleware(), controllers.CreatePhone)
			phoneRoutes.POST("/:phone_id/features", middleware.JWTAuthMiddleware(), controllers.CreateFeature)
			phoneRoutes.PUT("/:phone_id", middleware.JWTAuthMiddleware(), controllers.UpdatePhone)
			phoneRoutes.GET("/:phone_id", controllers.GetPhoneByID)
			phoneRoutes.DELETE("/:phone_id", middleware.JWTAuthMiddleware(), controllers.DeletePhone)
			phoneRoutes.PUT("/:phone_id/features/:feature_id", middleware.JWTAuthMiddleware(), controllers.UpdateFeature)
			phoneRoutes.DELETE("/:phone_id/features/:feature_id", middleware.JWTAuthMiddleware(), controllers.DeleteFeature)
		}

		reviewRoutes := api.Group("/reviews")

		{
			reviewRoutes.POST("/", middleware.JWTAuthMiddleware(), controllers.CreateReview)
			reviewRoutes.GET("/", controllers.GetAllReviews)
			reviewRoutes.GET("/:id", controllers.GetReviewByID)
			reviewRoutes.GET("/:phone_id", controllers.GetReviews)
			reviewRoutes.PUT("/:id", middleware.JWTAuthMiddleware(), controllers.UpdateReview)
			reviewRoutes.DELETE("/:id", middleware.JWTAuthMiddleware(), controllers.DeleteReview)
		}

		commentRoutes := api.Group("/comments")
		commentRoutes.Use(middleware.JWTAuthMiddleware())
		{
			commentRoutes.POST("/", controllers.CreateComment)
			commentRoutes.GET("/:review_id", controllers.GetComments)
			commentRoutes.PUT("/:id", controllers.UpdateComment)
			commentRoutes.DELETE("/:id", controllers.DeleteComment)
		}
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
