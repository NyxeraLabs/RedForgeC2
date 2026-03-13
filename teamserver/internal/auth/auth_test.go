package auth

import (
	"testing"
)

func TestGenerateAndValidateToken_Success(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, "admin", "admin", 10)
	if err != nil {
		t.Fatalf("expected token generation to succeed, got: %v", err)
	}

	claims, err := ValidateToken(secret, token)
	if err != nil {
		t.Fatalf("expected token validation to succeed, got: %v", err)
	}

	if claims.Username != "admin" {
		t.Fatalf("expected username %q, got %q", "admin", claims.Username)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected role %q, got %q", "admin", claims.Role)
	}
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, "admin", "admin", 10)
	if err != nil {
		t.Fatalf("expected token generation to succeed, got: %v", err)
	}

	_, err = ValidateToken("wrong-secret", token)
	if err == nil {
		t.Fatal("expected token validation to fail with wrong secret")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, "admin", "admin", -1)
	if err != nil {
		t.Fatalf("expected token generation to succeed, got: %v", err)
	}

	_, err = ValidateToken(secret, token)
	if err == nil {
		t.Fatal("expected token validation to fail for expired token")
	}
}
