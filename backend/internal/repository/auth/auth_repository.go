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
