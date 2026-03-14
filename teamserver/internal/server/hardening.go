package server

import (
	"crypto/subtle"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type hardeningConfig struct {
	allowedOrigins map[string]struct{}
	maxBodyBytes   int64
	loginRPM       int
}

func hardeningFromEnv() hardeningConfig {
	origins := make(map[string]struct{})
	raw := strings.TrimSpace(os.Getenv("REDFORGE_CORS_ORIGINS"))
	if raw == "" {
		// Secure-by-default for local development: allow only common Vite ports.
		raw = "http://localhost:5174,http://localhost:5173"
	}
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins[o] = struct{}{}
		}
	}

	maxBodyBytes := int64(25 << 20) // 25 MiB default to avoid accidental breakage
	if v := strings.TrimSpace(os.Getenv("REDFORGE_MAX_BODY_BYTES")); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed > 0 {
			maxBodyBytes = parsed
		}
	}

	loginRPM := 20
	if v := strings.TrimSpace(os.Getenv("REDFORGE_LOGIN_RPM")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			loginRPM = parsed
		}
	}

	return hardeningConfig{
		allowedOrigins: origins,
		maxBodyBytes:   maxBodyBytes,
		loginRPM:       loginRPM,
	}
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API-hardening headers: safe even if UI is served elsewhere.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		next.ServeHTTP(w, r)
	})
}

func withBodyLimit(maxBytes int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maxBytes > 0 && r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func withCORS(allowed map[string]struct{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && originAllowed(allowed, origin) {
			// If "*" is configured, return "*" (no credentials).
			if _, ok := allowed["*"]; ok {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func originAllowed(allowed map[string]struct{}, origin string) bool {
	if allowed == nil {
		return false
	}
	if _, ok := allowed["*"]; ok {
		return true
	}
	_, ok := allowed[origin]
	return ok
}

type ipRateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	rpm      int
	buckets  map[string]*ipBucket
	maxTrack int
}

type ipBucket struct {
	resetAt time.Time
	count   int
}

func newIPRateLimiter(rpm int, window time.Duration) *ipRateLimiter {
	if rpm <= 0 {
		rpm = 1
	}
	return &ipRateLimiter{
		window:   window,
		rpm:      rpm,
		buckets:  make(map[string]*ipBucket),
		maxTrack: 4096,
	}
}

func (l *ipRateLimiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.buckets) > l.maxTrack {
		// Very simple defense: if under memory pressure, drop oldest by nuking map.
		// This preserves availability while keeping brute-force protection best-effort.
		l.buckets = make(map[string]*ipBucket)
	}

	b, ok := l.buckets[ip]
	if !ok || now.After(b.resetAt) {
		l.buckets[ip] = &ipBucket{resetAt: now.Add(l.window), count: 1}
		return true
	}
	if b.count >= l.rpm {
		return false
	}
	b.count++
	return true
}

func withLoginRateLimit(l *ipRateLimiter, next http.Handler) http.Handler {
	if l == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		now := time.Now()
		if !l.allow(ip, now) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	// Avoid trusting X-Forwarded-For by default (deployment may not sanitize it).
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

// constantTimeEq is used for any security-sensitive comparisons we may add later.
func constantTimeEq(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

