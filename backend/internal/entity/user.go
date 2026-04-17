package entity

import "time"

type User struct {
	ID           int64      `db:"id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"` // Sekarang string biasa (bukan pointer) karena NOT NULL
	Phone        *string    `db:"phone"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
