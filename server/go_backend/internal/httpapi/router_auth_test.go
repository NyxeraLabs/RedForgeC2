package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"redforgec2/server/go_backend/internal/auth"
	"redforgec2/server/go_backend/internal/store"
)

func TestAuthProtectsOperatorEndpointsWhenEnabled(t *testing.T) {
	hash, err := auth.HashPassword("pw")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	authSvc, err := auth.New(true, "01234567890123456789012345678901", time.Hour, []auth.User{
		{Username: "admin", PasswordHash: hash, Role: auth.RoleAdmin},
	})
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}

	st := store.New()
	srv := httptest.NewServer(NewRouter(Deps{Store: st, Auth: authSvc}))
	t.Cleanup(srv.Close)

	// Unauthed list agents -> 401
	res, err := http.Get(srv.URL + "/api/v1/agents")
	if err != nil {
		t.Fatalf("get agents: %v", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.StatusCode)
	}

	// Login -> token
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "pw"})
	loginRes, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	defer loginRes.Body.Close()
	if loginRes.StatusCode != http.StatusOK {
		t.Fatalf("login status: %s", loginRes.Status)
	}
	var login map[string]any
	if err := json.NewDecoder(loginRes.Body).Decode(&login); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	token, _ := login["token"].(string)
	if token == "" {
		t.Fatalf("expected token")
	}

	// Authed list agents -> 200
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/agents", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	okRes, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("authed agents: %v", err)
	}
	okRes.Body.Close()
	if okRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", okRes.StatusCode)
	}
}

