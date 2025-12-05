package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token expired")
	ErrInvalidPassword = errors.New("invalid password")
)

// Claims represents the JWT claims
type Claims struct {
	UserName string `json:"userName,omitempty"`
	Room     string `json:"room,omitempty"`
	Identity string `json:"identity,omitempty"`
	UserType string `json:"userType,omitempty"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret string
}

// NewAuthService creates a new AuthService
func NewAuthService(jwtSecret string) *AuthService {
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

// CreateToken creates a new JWT token
func (s *AuthService) CreateToken(ctx context.Context, userName, room, identity, userType string, expiresIn time.Duration) (string, error) {
	if expiresIn == 0 {
		expiresIn = 24 * time.Hour // Default 24 hours
	}

	claims := &Claims{
		UserName: userName,
		Room:     room,
		Identity: identity,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// VerifyToken verifies a JWT token and returns the claims
func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// VerifyUser verifies a user by username and password
func (s *AuthService) VerifyUser(ctx context.Context, userName, password, expectedPassword string) error {
	if password != expectedPassword {
		return ErrInvalidPassword
	}
	return nil
}

// CreateRoomToken creates a token for room access
func (s *AuthService) CreateRoomToken(ctx context.Context, room string, expiresIn time.Duration) (string, error) {
	return s.CreateToken(ctx, "", room, "", "", expiresIn)
}
