package parentrouter

import (
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	parenthandler "github.com/ahmadzakyarifin/school-payment-system/internal/handler/parent"
	"github.com/ahmadzakyarifin/school-payment-system/internal/middleware"
	parentservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/parent"
	"github.com/gin-gonic/gin"
)

func Setup(a *app.App, rg *gin.RouterGroup) {
	service := parentservice.New()
	handler := parenthandler.New(service)

	parent := rg.Group("/parent")
	parent.Use(middleware.RequireRole("parent"))
	{
		parent.GET("/dashboard", handler.GetDashboard)
	}
}
