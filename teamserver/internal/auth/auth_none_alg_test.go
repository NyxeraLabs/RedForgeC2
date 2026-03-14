package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken_RejectsNoneAlgorithm(t *testing.T) {
	claims := &Claims{
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create a "none" token to ensure our ValidateToken enforces expected signing methods.
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("expected none-alg token creation to succeed, got: %v", err)
	}

	_, err = ValidateToken("any-secret", tokenStr)
	if err == nil {
		t.Fatal("expected token validation to fail for none-alg token")
	}
}

