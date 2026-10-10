package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Struct untuk klaim token
type JWTClaim struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Ambil secret dari ENV atau default
func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}
	return secret
}

// Generate token JWT dengan user_id dan role
func GenerateJWT(userID uint, role string) (string, error) {
	claims := JWTClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(getSecret()))
	if err != nil {
		return "", err
	}
	return t, nil
}

// Validasi token dan return (userID, role)
func ValidateToken(tokenString string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("metode token tidak valid")
			}
			return []byte(getSecret()), nil
		},
	)

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return 0, "", errors.New("token tidak valid")
	}

	return claims.UserID, claims.Role, nil
}
