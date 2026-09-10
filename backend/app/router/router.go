package router

import (
	"basic-app/auth"
	"basic-app/config"
	"basic-app/handler"
	mongorepo "basic-app/repository/mongo"
	"basic-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	mongo "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database, cfg config.Config) *gin.Engine {
	r := gin.Default()

	// Dependencies
	userRepository := mongorepo.NewUserRepository(database)
	authService := services.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	// Global/Public Endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":     true,
			"status": "application is up",
		})
	})

	// Authentication routes
	authRoutes := r.Group("/api/v1/auth")
	{
		// Public
		authRoutes.POST("/register", authHandler.Register)

		// Authenticated
		protected := authRoutes.Group("")
		protected.Use(auth.AuthMiddleware(
			cfg.JWTSecret,
			cfg.AuthAccessCookie,
		))

		protected.PATCH("/password", authHandler.ChangePassword)
	}

	return r
}
