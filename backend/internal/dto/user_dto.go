package dto

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string  `json:"name"     binding:"required,min=2"`
	Email    string  `json:"email"    binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	Phone    *string `json:"phone"    binding:"omitempty"`
	Role     string  `json:"role"     binding:"required,oneof=admin parent"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"     binding:"omitempty,min=2"`
	Email    *string `json:"email"    binding:"omitempty,email"`
	Phone    *string `json:"phone"    binding:"omitempty"`
	Role     *string `json:"role"     binding:"omitempty,oneof=admin parent"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}
