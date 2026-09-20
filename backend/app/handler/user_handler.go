package handler

import (
	"errors"
	"net/http"

	"basic-app/dto"
	"basic-app/repository"
	"basic-app/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(
	userService *services.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userIDString, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, err := h.userService.GetMe(
		c.Request.Context(),
		userIDString,
	)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to fetch user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userIDString, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req dto.UpdateUserProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.userService.UpdateProfile(
		c.Request.Context(),
		userIDString,
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

		case errors.Is(err, services.ErrEmailAlreadyExists),
			errors.Is(err, services.ErrAltEmailAlreadyExists),
			errors.Is(err, services.ErrUsernameAlreadyExists),
			errors.Is(err, services.ErrPhoneAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrAltEmailSameAsEmail),
			errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			// Validation errors from the service also reach here.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) CreateCustomer(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.userService.CreateCustomer(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists),
			errors.Is(err, services.ErrAltEmailAlreadyExists),
			errors.Is(err, services.ErrUsernameAlreadyExists),
			errors.Is(err, services.ErrPhoneAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrAltEmailSameAsEmail):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			// Validation errors reach here.
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "customer created successfully",
		"user":    user,
	})
}

func (h *UserHandler) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")

	var req dto.UpdateUserProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.userService.UpdateCustomer(
		c.Request.Context(),
		customerID,
		&req,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "customer not found",
			})

		case errors.Is(err, services.ErrInvalidUserRole):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "target user is not a customer",
			})

		case errors.Is(err, services.ErrEmailAlreadyExists),
			errors.Is(err, services.ErrAltEmailAlreadyExists),
			errors.Is(err, services.ErrUsernameAlreadyExists),
			errors.Is(err, services.ErrPhoneAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, services.ErrAltEmailSameAsEmail),
			errors.Is(err, services.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "customer updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) ListCustomers(c *gin.Context) {
	var query dto.ListCustomersQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid query parameters",
		})
		return
	}

	result, err := h.userService.ListCustomers(
		c.Request.Context(),
		query,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPage),
			errors.Is(err, services.ErrInvalidLimit):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to fetch customers",
			})
		}

		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) GetCustomerByID(c *gin.Context) {
	customerID := c.Param("id")

	user, err := h.userService.GetCustomerByID(
		c.Request.Context(),
		customerID,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "customer not found",
			})

		case errors.Is(err, services.ErrInvalidUserRole):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "target user is not a customer",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to fetch customer",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customer": user,
	})
}
