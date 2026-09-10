package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	userID string,
	role string,
	secret string,
	expiryHours int,
) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(expiryHours) * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(
	userID string,
	secret string,
	expiryDays int,
) (string, time.Time, error) {
	now := time.Now()

	expiresAt := now.Add(
		time.Duration(expiryDays) * 24 * time.Hour,
	)

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func ValidateToken(
	tokenString string,
	secret string,
) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
