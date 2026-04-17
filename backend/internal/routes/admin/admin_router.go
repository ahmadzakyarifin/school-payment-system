package adminrouter

import (
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	adminuserhandler "github.com/ahmadzakyarifin/school-payment-system/internal/handler/admin/user"
	"github.com/ahmadzakyarifin/school-payment-system/internal/middleware"
	adminuserrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/admin/user"
	adminuserservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/admin/user"
	"github.com/gin-gonic/gin"
)

func Setup(a *app.App, rg *gin.RouterGroup) {
	// Initialize User Stack (Modular)
	userRepo := adminuserrepo.New(a.DB)
	userService := adminuserservice.New(userRepo)
	userHandler := adminuserhandler.New(userService)

	admin := rg.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	{
		// User Management CRUD
		users := admin.Group("/users")
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/roles", userHandler.GetRoles)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PATCH("/:id", userHandler.UpdateUser)
			users.PATCH("/:id/status", userHandler.ToggleStatus)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
