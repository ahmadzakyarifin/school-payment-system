package authhandler

import (
	"errors"

	authdto "github.com/ahmadzakyarifin/school-payment-system/internal/dto/auth"
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
	var req authdto.LoginRequest

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

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req authdto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), req); err != nil {
		response.InternalServerError(c, "Gagal memproses permintaan reset password")
		return
	}

	response.OK(c, "Jika email terdaftar, instruksi reset password akan dikirim", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req authdto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Password berhasil diperbarui, silakan login kembali", nil)
}
