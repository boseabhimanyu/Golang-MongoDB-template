package dto

import "time"

// UpdateUserProfileRequest contains fields that can be
// partially updated.
//
// nil means the field was not provided.
type UpdateUserProfileRequest struct {
	FirstName    *string    `json:"firstName"`
	LastName     *string    `json:"lastName"`
	Username     *string    `json:"username"`
	Email        *string    `json:"email"`
	AltEmail     *string    `json:"altEmail"`
	Phone        *string    `json:"phone"`
	DateOfBirth  *time.Time `json:"dateOfBirth"`
	AddressLine1 *string    `json:"addressLine1"`
	AddressLine2 *string    `json:"addressLine2"`
	City         *string    `json:"city"`
	State        *string    `json:"state"`
	PinCode      *string    `json:"pinCode"`
}

type ChangeUserPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}
