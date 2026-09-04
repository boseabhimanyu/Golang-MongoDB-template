package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCustomer UserRole = "customer"
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	FirstName string `bson:"first_name" json:"firstName"`
	LastName  string `bson:"last_name" json:"lastName"`
	Username  string `bson:"username" json:"username"`
	Email     string `bson:"email" json:"email"`
	AltEmail  string `bson:"alt_email,omitempty" json:"altEmail,omitempty"`
	Phone     string `bson:"phone,omitempty" json:"phone,omitempty"`

	Role       UserRole `bson:"role" json:"role"`
	ProfilePic string   `bson:"profile_pic,omitempty" json:"profilePic,omitempty"`
	Status     bool     `bson:"status" json:"status"`

	PasswordHash     string `bson:"password_hash" json:"-"`
	RefreshTokenHash string `bson:"refresh_token_hash,omitempty" json:"-"`

	CreatedAt   time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updatedAt"`
	LastLoginAt *time.Time `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
}

func NewCustomerUser() User {
	now := time.Now().UTC()

	return User{
		Role:      RoleCustomer,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
