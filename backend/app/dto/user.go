package dto

import "basic-app/models"

// UpdateUserProfileRequest contains fields that can be
// partially updated.
//
// nil means the field was not provided.
type UpdateUserProfileRequest struct {
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Username     *string `json:"username"`
	Email        *string `json:"email"`
	AltEmail     *string `json:"altEmail"`
	Phone        *string `json:"phone"`
	DateOfBirth  *string `json:"dateOfBirth"`
	AddressLine1 *string `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PinCode      *string `json:"pinCode"`
}

type ChangeUserPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

type ListCustomersQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Status *bool  `form:"status"`
	Search string `form:"search"`
}

type CustomerListResponse struct {
	Customers  []models.User `json:"customers"`
	Pagination Pagination    `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

type UserStatusRequest struct {
	Status *bool `json:"status"`
}
