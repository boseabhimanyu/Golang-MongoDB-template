package handler

import (
	"errors"
	"net/http"

	"basic-app/dto"
	"basic-app/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(
	authService *services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles public customer registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.authService.RegisterUser(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":        user.ID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"username":  user.Username,
			"email":     user.Email,
			"altEmail":  user.AltEmail,
			"phone":     user.Phone,
			"role":      user.Role,
			"status":    user.Status,
		},
	})
}

// ChangePassword handles a user's own password change.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.authService.ChangePassword(
		c.Request.Context(),
		userID,
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "current password is incorrect",
			})

		case errors.Is(err, services.ErrInvalidUserID):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password changed successfully",
	})
}
