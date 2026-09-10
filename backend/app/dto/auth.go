package dto

// RegisterRequest contains the allowed fields
// for public customer registration.
type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Username string `json:"username"`

	Email    string `json:"email"`
	AltEmail string `json:"altEmail"`
	Phone    string `json:"phone"`

	Password string `json:"password"`
}

// ChangePasswordRequest contains the current password
// and the new password for a password change.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
