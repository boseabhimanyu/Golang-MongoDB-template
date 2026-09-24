package router

import (
	"basic-app/auth"
	"basic-app/config"
	"basic-app/handler"
	mongorepo "basic-app/repository/mongo"
	"basic-app/services"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mongo "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database, cfg config.Config) *gin.Engine {
	r := gin.Default()

	r.Static("/Uploads", "./Uploads") // Expose file uploads

	// Dependencies
	userRepository := mongorepo.NewUserRepository(database)
	authService := services.NewAuthService(userRepository, cfg)
	refreshToken := func(
		ctx context.Context,
		refreshToken string,
	) (string, string, time.Time, error) {
		result, err := authService.RefreshToken(ctx, refreshToken)
		if err != nil {
			return "", "", time.Time{}, err
		}

		return result.AccessToken,
			result.RefreshToken,
			result.RefreshExpiry,
			nil
	}

	authMiddleware := auth.AuthMiddleware(
		cfg.JWTSecret,
		cfg.AuthAccessCookie,
		cfg.AuthRefreshCookie,
		cfg.CookieSecure,
		cfg.JWTExpiryHours,
		refreshToken,
	)
	userService := services.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)

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
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)

		// Authenticated
		protected := authRoutes.Group("")
		protected.Use(authMiddleware)

		protected.PATCH("/password", authHandler.ChangePassword)
		//	protected.POST("/refresh", authHandler.Refresh)
		protected.POST("/logout", authHandler.Logout)
		protected.PATCH("/me", userHandler.UpdateProfile)
		protected.PATCH("/me/image", userHandler.UpdateProfilePic)
		protected.GET("/me", userHandler.Me)
	}

	customerRoutes := r.Group("/api/v1/customers")

	customerRoutes.Use(
		authMiddleware,
		auth.RequireRoles("admin"),
	)

	customerRoutes.POST("", userHandler.CreateCustomer)

	//--------------------------------------------------------------
	customerRoutes.GET("", userHandler.ListCustomers)

	// List customers.
	//
	// Pagination:
	// GET /api/v1/customers?page=1&limit=20
	//
	// Filter by account status:
	// GET /api/v1/customers?status=true
	// GET /api/v1/customers?status=false
	//
	// Search across first name, last name, username, email,
	// alternate email, and phone:
	// GET /api/v1/customers?search=rahul
	//
	// Filters can be combined:
	// GET /api/v1/customers?page=1&limit=20&status=true&search=rahul

	//--------------------------------------------------------------

	customerRoutes.GET("/:id", userHandler.GetCustomerByID)
	customerRoutes.PATCH("/:id", userHandler.UpdateCustomer)
	customerRoutes.PATCH("/:id/status", userHandler.UpdateUserStatus)
	customerRoutes.PATCH("/:id/password", userHandler.ChangeUserPassword)

	return r
}
