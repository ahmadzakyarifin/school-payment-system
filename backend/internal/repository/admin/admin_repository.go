package adminrepo

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/ahmadzakyarifin/school-payment-system/internal/dto"
	"github.com/ahmadzakyarifin/school-payment-system/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AdminRepository interface {
	FindAll(ctx context.Context, req dto.UserListRequest) ([]entity.User, int64, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	GetRoles(ctx context.Context) ([]string, error)
	Create(ctx context.Context, user *entity.User) (int64, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateStatus(ctx context.Context, id int64, isActive bool) error
	Delete(ctx context.Context, id int64) error
	CheckDuplicate(ctx context.Context, email string, excludeID int64) (bool, error)
}

type adminRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindAll(ctx context.Context, req dto.UserListRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	where := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if req.Search != "" {
		// Sekarang fokus pencarian hanya di name dan email (karena username sudah dihapus)
		where += " AND MATCH(name, email) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, req.Search+"*")
	}

	if req.Role != "" {
		where += " AND role = ?"
		args = append(args, req.Role)
	}

	if req.IsActive != nil {
		where += " AND is_active = ?"
		args = append(args, *req.IsActive)
	}

	countQuery := "SELECT COUNT(*) FROM users " + where
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf("SELECT * FROM users %s ORDER BY created_at DESC LIMIT %d OFFSET %d", where, req.Limit, offset)

	err = r.db.SelectContext(ctx, &users, query, args...)
	return users, total, err
}

func (r *adminRepository) CheckDuplicate(ctx context.Context, email string, excludeID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ? AND deleted_at IS NULL)`
	err := r.db.GetContext(ctx, &exists, query, email, excludeID)
	return exists, err
}

func (r *adminRepository) GetRoles(ctx context.Context) ([]string, error) {
	var columnType string
	query := `SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='users' AND COLUMN_NAME='role' AND TABLE_SCHEMA=(SELECT DATABASE())`
	err := r.db.GetContext(ctx, &columnType, query)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`'([^']*)'`)
	matches := re.FindAllStringSubmatch(columnType, -1)
	var roles []string
	for _, m := range matches {
		roles = append(roles, m[1])
	}
	return roles, nil
}

func (r *adminRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL"
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *adminRepository) Create(ctx context.Context, user *entity.User) (int64, error) {
	query := `INSERT INTO users (name, email, password_hash, phone, role, is_active) VALUES (:name, :email, :password_hash, :phone, :role, :is_active)`
	res, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *adminRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name=:name, email=:email, phone=:phone, role=:role, is_active=:is_active WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *adminRepository) UpdateStatus(ctx context.Context, id int64, isActive bool) error {
	query := "UPDATE users SET is_active = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, isActive, id)
	return err
}

func (r *adminRepository) Delete(ctx context.Context, id int64) error {
	query := "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
