package authhandler

import (
	"errors"

	"github.com/ahmadzakyarifin/school-payment-system/internal/dto"
	authservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/auth"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service authservice.AuthService
}

func New(service authservice.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) ||
			errors.Is(err, authservice.ErrAccountInactive) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalServerError(c, "Terjadi kesalahan pada server")
		return
	}

	response.OK(c, "Login berhasil", result)
}
