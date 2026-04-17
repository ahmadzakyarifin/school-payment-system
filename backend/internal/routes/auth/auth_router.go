package authrouter

import (
	"os"

	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	authhandler "github.com/ahmadzakyarifin/school-payment-system/internal/handler/auth"
	authrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/auth"
	authservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/auth"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/mailer"
	"github.com/gin-gonic/gin"
)

func Setup(a *app.App, rg *gin.RouterGroup) {
	resendMailer := mailer.NewResend(
		os.Getenv("RESEND_API_KEY"),
		os.Getenv("MAIL_FROM"),
	)

	repo := authrepo.New(a.DB)
	service := authservice.New(repo, resendMailer, os.Getenv("JWT_SECRET"))
	handler := authhandler.New(service)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/reset-password", handler.ResetPassword)
	}
}
