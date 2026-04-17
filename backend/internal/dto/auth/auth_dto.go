package authdto

import userdto "github.com/ahmadzakyarifin/school-payment-system/internal/dto/user"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string               `json:"token"`
	User  userdto.UserResponse `json:"user"`
}
