package auth

import (
	"testing"
	"time"
)

func TestLoginAndParse(t *testing.T) {
	hash, err := HashPassword("pass1234")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	svc, err := New(true, "01234567890123456789012345678901", time.Hour, []User{
		{Username: "alice", PasswordHash: hash, Role: RoleOperator},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	token, exp, role, err := svc.Login("alice", "pass1234")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" || exp.IsZero() || role != RoleOperator {
		t.Fatalf("unexpected login result")
	}
	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.Username != "alice" || claims.Role != RoleOperator {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestRejectsShortSecret(t *testing.T) {
	_, err := New(true, "short", time.Hour, []User{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

