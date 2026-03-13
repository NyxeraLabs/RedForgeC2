package users

import (
	"context"
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
