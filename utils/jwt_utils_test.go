package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-enough-entropy")
	token, err := GenerateJWT(42, "admin")
	if err != nil {
		t.Fatal(err)
	}
	userID, role, err := ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if userID != 42 || role != "admin" {
		t.Fatalf("got userID=%d role=%q", userID, role)
	}
}

func TestValidateTokenRejectsInvalidTokens(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-enough-entropy")

	t.Run("tampered", func(t *testing.T) {
		token, err := GenerateJWT(1, "user")
		if err != nil {
			t.Fatal(err)
		}
		if _, _, err := ValidateToken(token + "x"); err == nil {
			t.Fatal("tampered token accepted")
		}
	})

	t.Run("expired", func(t *testing.T) {
		claims := JWTClaim{
			UserID: 1,
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			t.Fatal(err)
		}
		if _, _, err := ValidateToken(token); err == nil {
			t.Fatal("expired token accepted")
		}
	})
}
