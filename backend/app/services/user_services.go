package services

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"basic-app/dto"
	"basic-app/models"
	"basic-app/repository"
	"basic-app/validation"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidUserRole       = errors.New("invalid user role")
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	return s.userRepository.FindByID(
		ctx,
		userID,
	)
}

// GetByEmail returns a user by email.
func (s *UserService) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	email = strings.ToLower(
		strings.TrimSpace(email),
	)

	return s.userRepository.FindByEmail(
		ctx,
		email,
	)

}

// GetByUsername returns a user by username.
func (s *UserService) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	username = strings.TrimSpace(
		username,
	)

	return s.userRepository.FindByUsername(
		ctx,
		username,
	)

}

// UpdateProfile updates general user information.
func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID string,
	req *dto.UpdateUserProfileRequest,
) error {
	currentUser, err := s.userRepository.FindByID(
		ctx,
		userID,
	)
	if err != nil {
		return err
	}

	// Validate first name.
	firstName, err := validation.ValidateName(
		*req.FirstName,
		"first name",
	)
	if err != nil {
		return err
	}

	// Validate last name.
	lastName, err := validation.ValidateName(
		*req.LastName,
		"last name",
	)
	if err != nil {
		return err
	}

	// Validate username.
	username, err := validation.ValidateUsername(
		*req.Username,
	)
	if err != nil {
		return err
	}

	// Validate required email.
	email, err := validation.ValidateEmail(
		*req.Email,
		true,
	)
	if err != nil {
		return err
	}

	// Validate optional alternate email.
	altEmail, err := validation.ValidateEmail(
		*req.AltEmail,
		false,
	)
	if err != nil {
		return err
	}

	// Validate phone.
	phone, err := validation.ValidatePhone(
		*req.Phone,
	)
	if err != nil {
		return err
	}

	// Check whether another user already uses this email.
	if email != currentUser.Email {
		existingUser, err := s.userRepository.FindByEmail(
			ctx,
			email,
		)

		if err == nil &&
			existingUser != nil &&
			existingUser.ID != currentUser.ID {
			return ErrEmailAlreadyExists
		}

		if err != nil &&
			!errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
	}

	// Check whether another user already uses this username.
	if username != currentUser.Username {
		existingUser, err := s.userRepository.FindByUsername(
			ctx,
			username,
		)

		if err == nil &&
			existingUser != nil &&
			existingUser.ID != currentUser.ID {
			return ErrUsernameAlreadyExists
		}

		if err != nil &&
			!errors.Is(err, repository.ErrUserNotFound) {
			return err
		}
	}

	// Update only allowed fields.
	currentUser.FirstName = firstName
	currentUser.LastName = lastName
	currentUser.Username = username

	currentUser.Email = email
	currentUser.AltEmail = altEmail
	currentUser.Phone = phone

	currentUser.DateOfBirth = req.DateOfBirth

	currentUser.AddressLine1 = *req.AddressLine1
	currentUser.AddressLine2 = *req.AddressLine2
	currentUser.City = *req.City
	currentUser.State = *req.State
	currentUser.PinCode = *req.PinCode

	return s.userRepository.Update(
		ctx,
		currentUser,
	)

}

// UpdateProfilePic updates only the profile image URL/path.
func (s *UserService) UpdateProfilePic(
	ctx context.Context,
	userID string,
	fileHeader *multipart.FileHeader,
) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidUserID
	}

	if fileHeader == nil {
		return errors.New("profile picture is required")
	}

	const maxFileSize = 5 * 1024 * 1024

	if fileHeader.Size > maxFileSize {
		return errors.New("profile picture must not exceed 5 MB")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		// allowed
	default:
		return errors.New("profile picture must be JPEG, PNG, or WebP")
	}

	// Make sure the target user exists.
	user, err := s.userRepository.FindByID(ctx, strings.TrimSpace(userID))
	if err != nil {
		return err
	}

	_ = user

	uploadDir := filepath.Join("Uploads", "profiles")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	fileID := uuid.New().String()
	filename := fileID + ext

	filePath := filepath.Join(uploadDir, filename)

	if err := saveUploadedFile(fileHeader, filePath); err != nil {
		return err
	}

	return s.userRepository.UpdateProfilePic(
		ctx,
		userID,
		filePath,
	)
}

func saveUploadedFile(
	fileHeader *multipart.FileHeader,
	destination string,
) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// DeactivateUser disables a user account.
func (s *UserService) DeactivateUser(
	ctx context.Context,
	userID string,
) error {
	return s.userRepository.Deactivate(
		ctx,
		userID,
	)
}

// DeleteUser permanently deletes a user.
func (s *UserService) DeleteUser(
	ctx context.Context,
	userID string,
) error {
	return s.userRepository.Delete(
		ctx,
		userID,
	)
}

func (s *UserService) CreateCustomer(
	ctx context.Context,
	req *dto.RegisterRequest,
) (*models.User, error) {

	if req == nil {
		return nil, errors.New("create customer request is required")
	}

	// Validate required fields.
	firstName, err := validation.ValidateName(req.FirstName, "first name")
	if err != nil {
		return nil, err
	}

	lastName, err := validation.ValidateName(req.LastName, "last name")
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

	phone, err := validation.ValidatePhone(req.Phone)
	if err != nil {
		return nil, err
	}

	if err := validation.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Validate optional alternate email.
	altEmail, err := validation.ValidateEmail(req.AltEmail, false)
	if err != nil {
		return nil, err
	}

	// Check duplicate email.
	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Check duplicate username.
	existingUser, err = s.userRepository.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUsernameAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Hash password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	// Always create a customer.
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

// ChangeUserPassword allows an administrator to reset the password
// of a customer user.
func (s *UserService) ChangeUserPassword(
	ctx context.Context,
	userID string,
	req *dto.ChangeUserPasswordRequest,
) error {

	if req == nil {
		return errors.New("change user password request is required")
	}

	if err := validation.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidUserID
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// For now, administrators can only reset customer passwords.
	if user.Role != models.RoleCustomer {
		return ErrInvalidUserRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
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
