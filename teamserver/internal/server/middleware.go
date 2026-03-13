package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
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
