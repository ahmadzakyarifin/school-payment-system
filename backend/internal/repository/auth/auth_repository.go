package authrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ahmadzakyarifin/school-payment-system/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	SaveResetToken(ctx context.Context, email, token string) error
	GetResetToken(ctx context.Context, token string) (string, error)
	DeleteResetToken(ctx context.Context, token string) error
	UpdatePassword(ctx context.Context, email, hashedPassword string) error
}

type authRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, name, email, phone, password_hash, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("authrepo.FindByEmail: %w", err)
	}

	return &user, nil
}

func (r *authRepository) SaveResetToken(ctx context.Context, email, token string) error {
	// Hapus token lama jika ada sebelum buat baru
	r.db.ExecContext(ctx, "DELETE FROM password_resets WHERE email = ?", email)

	query := "INSERT INTO password_resets (email, token, expires_at) VALUES (?, ?, DATE_ADD(NOW(), INTERVAL 15 MINUTE))"
	_, err := r.db.ExecContext(ctx, query, email, token)
	return err
}

func (r *authRepository) GetResetToken(ctx context.Context, token string) (string, error) {
	var email string
	query := "SELECT email FROM password_resets WHERE token = ? AND expires_at > NOW()"
	err := r.db.GetContext(ctx, &email, query, token)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return email, err
}

func (r *authRepository) DeleteResetToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM password_resets WHERE token = ?", token)
	return err
}

func (r *authRepository) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET password_hash = ? WHERE email = ?", hashedPassword, email)
	return err
}
