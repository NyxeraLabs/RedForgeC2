package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/users"
)

// AuthMiddleware validates JWT tokens and populates the request context.
func AuthMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Expect "Bearer <token>"
		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil || token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(secret, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddlewareWithAPIToken works like AuthMiddleware but will also accept API tokens stored in the database.
func AuthMiddlewareWithAPIToken(secret string, store *users.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// Browsers cannot set custom headers during the WebSocket handshake or file downloads.
		// For /api/ws and /api/operator/files/*/download, accept `?token=` as an alternative to the Authorization header.
		if authHeader == "" && (r.URL.Path == "/api/ws" || (strings.Contains(r.URL.Path, "/api/operator/files/") && strings.Contains(r.URL.Path, "/download"))) {
			if q := r.URL.Query().Get("token"); q != "" {
				authHeader = "Bearer " + q
			}
		}
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Expect "Bearer <token>"
		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil || token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(secret, token)
		if err != nil {
			user, err2 := store.ValidateAPIToken(r.Context(), token)
			if err2 != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims = &auth.Claims{Username: user.Username, Role: string(user.Role)}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole ensures the JWT contains the required role.
func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(contextKey("claims"))
		claims, ok := val.(*auth.Claims)
		if !ok || claims.Role != role {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAnyRole(roles []string, next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(contextKey("claims"))
		claims, ok := val.(*auth.Claims)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if _, ok := allowed[claims.Role]; !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ClaimsFromRequest(r *http.Request) (*auth.Claims, bool) {
	val := r.Context().Value(contextKey("claims"))
	claims, ok := val.(*auth.Claims)
	return claims, ok
}

// contextKey is used for storing values in request context.
type contextKey string

// auditResponseWriter wraps http.ResponseWriter to capture status and bytes written.
type auditResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (w *auditResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

// withAuditLog logs each HTTP request/response cycle for auditing purposes.
func withAuditLog(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		awr := &auditResponseWriter{ResponseWriter: w}

		next.ServeHTTP(awr, r)

		if awr.status == 0 {
			awr.status = http.StatusOK
		}

		user := "anonymous"
		if claims, ok := ClaimsFromRequest(r); ok {
			user = claims.Username
		}

		logger.Printf(
			"audit user=%s ip=%s method=%s path=%s status=%d bytes=%d duration_ms=%d",
			user,
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			awr.status,
			awr.bytesWritten,
			time.Since(start).Milliseconds(),
		)
	})
}
