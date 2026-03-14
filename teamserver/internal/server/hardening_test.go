package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithLoginRateLimit_BlocksAfterLimit(t *testing.T) {
	lim := newIPRateLimiter(1, time.Minute)

	calls := 0
	h := withLoginRateLimit(lim, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "203.0.113.10:1234"

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}

	if calls != 1 {
		t.Fatalf("expected handler to be called %d time, got %d", 1, calls)
	}
}

func TestWithCORS_DisallowedOrigin_NoAllowOriginHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://evil.example")
	w := httptest.NewRecorder()

	allowed := map[string]struct{}{"http://localhost:5174": {}}
	h := withCORS(allowed, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no ACAO header for disallowed origin, got %q", got)
	}
}

