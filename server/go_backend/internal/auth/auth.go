package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleOperator Role = "operator"
	RoleAnalyst Role = "analyst"
	RoleViewer  Role = "viewer"
)

type User struct {
	Username     string
	PasswordHash []byte
	Role         Role
}

type Service struct {
	enabled  bool
	secret   []byte
	ttl      time.Duration
	userByName map[string]User
}

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

func New(enabled bool, jwtSecret string, ttl time.Duration, users []User) (*Service, error) {
	s := &Service{
		enabled:    enabled,
		secret:     []byte(jwtSecret),
		ttl:        ttl,
		userByName: map[string]User{},
	}
	for _, u := range users {
		if u.Username == "" || len(u.PasswordHash) == 0 || u.Role == "" {
			return nil, errors.New("invalid user config")
		}
		s.userByName[u.Username] = u
	}
	if s.enabled && len(s.secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 bytes when auth is enabled")
	}
	if s.ttl <= 0 {
		s.ttl = 8 * time.Hour
	}
	return s, nil
}

func (s *Service) Enabled() bool { return s.enabled }

func (s *Service) Login(username, password string) (token string, expiresAt time.Time, role Role, err error) {
	if !s.enabled {
		return "", time.Time{}, "", errors.New("auth disabled")
	}
	u, ok := s.userByName[username]
	if !ok {
		return "", time.Time{}, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return "", time.Time{}, "", errors.New("invalid credentials")
	}

	now := time.Now().UTC()
	exp := now.Add(s.ttl)
	claims := Claims{
		Username: username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, "", err
	}
	return signed, exp, u.Role, nil
}

func (s *Service) Parse(token string) (*Claims, error) {
	if !s.enabled {
		return nil, errors.New("auth disabled")
	}
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func HashPassword(pw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
}

