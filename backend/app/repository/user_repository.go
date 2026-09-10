package repository

import (
	"context"
	"errors"

	"basic-app/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error

	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)

	Update(ctx context.Context, user *models.User) error

	UpdateProfilePic(
		ctx context.Context,
		userID string,
		profilePic string,
	) error

	UpdatePassword(
		ctx context.Context,
		userID string,
		passwordHash string,
	) error

	UpdateRefreshToken(
		ctx context.Context,
		userID string,
		refreshTokenHash string,
	) error

	Deactivate(ctx context.Context, userID string) error

	Delete(ctx context.Context, userID string) error
}
