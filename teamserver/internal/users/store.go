package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleObserver Role = "observer"
)

type User struct {
	Username     string    `json:"username"`
	Role         Role      `json:"role"`
	DisplayName  string    `json:"display_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PasswordHash string    `json:"-"`
}

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) EnsureAdmin(ctx context.Context, username string, password string, reset bool) error {
	if username == "" || password == "" {
		return errors.New("admin username/password required")
	}

	_, existingErr := s.GetByUsername(ctx, username)
	if existingErr != nil && !errors.Is(existingErr, pgx.ErrNoRows) {
		return existingErr
	}
	if existingErr == nil && !reset {
		return nil
	}

	hash, hashErr := HashPassword(password)
	if hashErr != nil {
		return hashErr
	}

	if existingErr == nil {
		_, execErr := s.pool.Exec(ctx, `
			UPDATE users
			SET password_hash=$2, role=$3, updated_at=NOW()
			WHERE username=$1
		`, username, hash, RoleAdmin)
		return execErr
	}

	_, execErr := s.pool.Exec(ctx, `
		INSERT INTO users (username, password_hash, role, display_name)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (username) DO NOTHING
	`, username, hash, RoleAdmin, "Administrator")
	return execErr
}

func HashPassword(password string) (string, error) {
	if len(password) < 10 {
		return "", fmt.Errorf("password too short (min 10 chars)")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Store) Create(ctx context.Context, username string, password string, role Role, displayName string) error {
	if username == "" {
		return errors.New("username required")
	}
	if role != RoleAdmin && role != RoleOperator && role != RoleObserver {
		return fmt.Errorf("invalid role: %s", role)
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO users (username, password_hash, role, display_name)
		VALUES ($1,$2,$3,$4)
	`, username, hash, role, displayName)
	return err
}

func (s *Store) List(ctx context.Context) ([]User, error) {
	rows, err := s.pool.Query(ctx, `SELECT username, role, display_name, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Role, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt); err != nil {
			continue
		}
		out = append(out, u)
	}
	return out, nil
}

func (s *Store) GetByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT username, password_hash, role, display_name, created_at, updated_at
		FROM users
		WHERE username=$1
	`, username).Scan(&u.Username, &u.PasswordHash, &u.Role, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) Authenticate(ctx context.Context, username string, password string) (*User, error) {
	u, err := s.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *Store) UpdateProfile(ctx context.Context, username string, displayName string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE users
		SET display_name=$2, updated_at=NOW()
		WHERE username=$1
	`, username, displayName)
	return err
}

func (s *Store) ChangePassword(ctx context.Context, username string, oldPassword string, newPassword string) error {
	u, err := s.Authenticate(ctx, username, oldPassword)
	if err != nil {
		return err
	}
	_ = u
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		UPDATE users
		SET password_hash=$2, updated_at=NOW()
		WHERE username=$1
	`, username, hash)
	return err
}

// CreateAPIToken generates a new API token for the given user and stores a hash in the database.
// The returned token is the only copy; callers must store it securely.
func (s *Store) CreateAPIToken(ctx context.Context, username string, description string, expiresAt *time.Time) (string, error) {
	// Ensure user exists.
	_, err := s.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	// Use URL-safe base64 without padding for token transport.
	token := base64.RawURLEncoding.EncodeToString(raw)
	// Store only a hash of the token.
	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])

	_, err = s.pool.Exec(ctx, `
		INSERT INTO api_tokens (token_hash, username, description, expires_at)
		VALUES ($1,$2,$3,$4)
	`, hashHex, username, description, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateAPIToken verifies an API token and returns the owning user if valid.
type APIToken struct {
	TokenHash   string     `json:"token_hash"`
	Username    string     `json:"username"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsed    *time.Time `json:"last_used,omitempty"`
}

func (s *Store) ValidateAPIToken(ctx context.Context, token string) (*User, error) {
	if token == "" {
		return nil, errors.New("token required")
	}
	// Token is stored as SHA256 of the raw token.
	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])

	var u User
	var expiresAt *time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT u.username, u.role, u.display_name, u.created_at, u.updated_at, t.expires_at
		FROM api_tokens t
		JOIN users u ON u.username = t.username
		WHERE t.token_hash = $1
	`, hashHex).Scan(&u.Username, &u.Role, &u.DisplayName, &u.CreatedAt, &u.UpdatedAt, &expiresAt)
	if err != nil {
		return nil, err
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		return nil, errors.New("token expired")
	}

	// Update last used timestamp (best-effort).
	_, _ = s.pool.Exec(ctx, `
		UPDATE api_tokens SET last_used = NOW() WHERE token_hash = $1
	`, hashHex)

	return &u, nil
}

func (s *Store) ListAPITokens(ctx context.Context) ([]APIToken, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT token_hash, username, description, created_at, expires_at, last_used
		FROM api_tokens
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]APIToken, 0)
	for rows.Next() {
		var t APIToken
		if err := rows.Scan(&t.TokenHash, &t.Username, &t.Description, &t.CreatedAt, &t.ExpiresAt, &t.LastUsed); err != nil {
			continue
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (s *Store) RevokeAPIToken(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM api_tokens WHERE token_hash = $1`, tokenHash)
	return err
}
