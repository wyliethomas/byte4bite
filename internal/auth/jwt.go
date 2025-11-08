package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	PantryID *uuid.UUID `json:"pantry_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey   string
	expiryHours int
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, expiryHours int) *JWTService {
	return &JWTService{
		secretKey:   secretKey,
		expiryHours: expiryHours,
	}
}

// GenerateToken generates a new JWT token
func (s *JWTService) GenerateToken(userID uuid.UUID, email, role string, pantryID *uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.expiryHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		PantryID: pantryID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiry
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims but new expiry
	return s.GenerateToken(claims.UserID, claims.Email, claims.Role, claims.PantryID)
}
