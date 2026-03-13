package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/operator/agents", nil)
	w := httptest.NewRecorder()

	handler := AuthMiddleware("secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	credsToken, _ := auth.GenerateToken("secret", "admin", "admin", 10)

	req := httptest.NewRequest(http.MethodGet, "/api/operator/agents", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()

	handler := AuthMiddleware("secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Ensure valid token works
	req = httptest.NewRequest(http.MethodGet, "/api/operator/agents", nil)
	req.Header.Set("Authorization", "Bearer "+credsToken)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	// create a request with claims but wrong role
	req := httptest.NewRequest(http.MethodGet, "/api/operator/agents", nil)
	ctx := req.Context()
	claims := &auth.Claims{Username: "user", Role: "user"}
	ctx = contextWithClaims(ctx, claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler := RequireRole("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/operator/agents", nil)
	ctx := req.Context()
	claims := &auth.Claims{Username: "admin", Role: "admin"}
	ctx = contextWithClaims(ctx, claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler := RequireRole("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// contextWithClaims is a helper to mimic AuthMiddleware behavior for role tests.
func contextWithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, contextKey("claims"), claims)
}
