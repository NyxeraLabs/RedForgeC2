package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz_OK(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	s.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", ct)
	}
	if body := w.Body.String(); body != `{"status":"ok"}` {
		t.Fatalf("expected body %q, got %q", `{"status":"ok"}`, body)
	}
}

func TestWithCORS_Preflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/anything", nil)
	req.Header.Set("Origin", "http://localhost:5174")
	w := httptest.NewRecorder()

	allowed := map[string]struct{}{"http://localhost:5174": {}}
	h := withCORS(allowed, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("expected handler not to be called for OPTIONS preflight")
	}))
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5174" {
		t.Fatalf("expected ACAO %q, got %q", "http://localhost:5174", got)
	}
}
