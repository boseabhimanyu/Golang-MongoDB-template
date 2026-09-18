package services

import (
	"basic-app/dto"
	"basic-app/models"
	"basic-app/repository"
	"basic-app/validation"
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(
	userRepository repository.UserRepository,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

// RegisterCustomer creates a new customer account.

func (s *AuthService) RegisterUser(
	ctx context.Context,
	req *dto.RegisterRequest,
) (*models.User, error) {
	firstName, err := validation.ValidateName(
		req.FirstName,
		"first name",
	)
	if err != nil {
		return nil, err
	}

	lastName, err := validation.ValidateName(
		req.LastName,
		"last name",
	)
	if err != nil {
		return nil, err
	}

	username, err := validation.ValidateUsername(req.Username)
	if err != nil {
		return nil, err
	}

	email, err := validation.ValidateEmail(req.Email, true)
	if err != nil {
		return nil, err
	}

	altEmail, err := validation.ValidateEmail(req.AltEmail, false)
	if err != nil {
		return nil, err
	}

	if altEmail != "" && altEmail == email {
		return nil, ErrAltEmailSameAsEmail
	}

	phone, err := validation.ValidatePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	if err := validation.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Checks primary email against both email and alternate-email fields.
	existingUser, err := s.userRepository.FindByAnyEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Checks optional alternate email against both email fields.
	if altEmail != "" {
		existingUser, err = s.userRepository.FindByAnyEmail(
			ctx,
			altEmail,
		)
		if err == nil && existingUser != nil {
			return nil, ErrAltEmailAlreadyExists
		}

		if err != nil && !errors.Is(
			err,
			repository.ErrUserNotFound,
		) {
			return nil, err
		}
	}

	existingUser, err = s.userRepository.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUsernameAlreadyExists
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	existingUser, err = s.userRepository.FindByPhone(ctx, phone)
	if err == nil && existingUser != nil {
		return nil, ErrPhoneAlreadyExists
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := models.NewCustomerUser()

	user.FirstName = firstName
	user.LastName = lastName
	user.Username = username
	user.Email = email
	user.AltEmail = altEmail
	user.Phone = phone
	user.PasswordHash = string(passwordHash)

	if err := s.userRepository.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePassword changes a user's password.
func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID string,
	req *dto.ChangePasswordRequest,
) error {
	if req == nil {
		return errors.New("change password request is required")
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidUserID
	}

	// Do not trim passwords.
	currentPassword := req.CurrentPassword
	newPassword := req.NewPassword

	if err := validation.ValidatePassword(newPassword); err != nil {
		return err
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify the user's current password.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(currentPassword),
	); err != nil {
		return ErrInvalidPassword
	}

	// Hash the new password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.userRepository.UpdatePassword(
		ctx,
		userID,
		string(passwordHash),
	)
}

// UpdateRefreshToken stores the hashed refresh token.
func (s *AuthService) UpdateRefreshToken(
	ctx context.Context,
	userID string,
	refreshTokenHash string,
) error {
	return s.userRepository.UpdateRefreshToken(
		ctx,
		userID,
		refreshTokenHash,
	)
}

// ClearRefreshToken removes the current refresh token.
func (s *AuthService) ClearRefreshToken(
	ctx context.Context,
	userID string,
) error {
	return s.userRepository.UpdateRefreshToken(
		ctx,
		userID,
		"",
	)
}
