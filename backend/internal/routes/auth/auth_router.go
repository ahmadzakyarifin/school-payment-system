package authrouter

import (
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	authhandler "github.com/ahmadzakyarifin/school-payment-system/internal/handler/auth"
	authrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/auth"
	authservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/auth"
	"github.com/gin-gonic/gin"
)

// Setup mendaftarkan semua route AUTH ke dalam router group.
// Semua route di sini bersifat PUBLIC — tidak memerlukan JWT.
func Setup(a *app.App, rg *gin.RouterGroup) {
	repo := authrepo.New(a.DB)
	service := authservice.New(repo, a.Config.JWTSecret)
	handler := authhandler.New(service)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
	}
}
