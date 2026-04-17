package adminrouter

import (
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	adminhandler "github.com/ahmadzakyarifin/school-payment-system/internal/handler/admin"
	"github.com/ahmadzakyarifin/school-payment-system/internal/middleware"
	adminrepo "github.com/ahmadzakyarifin/school-payment-system/internal/repository/admin"
	adminservice "github.com/ahmadzakyarifin/school-payment-system/internal/service/admin"
	"github.com/gin-gonic/gin"
)

func Setup(a *app.App, rg *gin.RouterGroup) {
	repo := adminrepo.New(a.DB)
	service := adminservice.New(repo)
	handler := adminhandler.New(service)

	admin := rg.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	{
		// User Management CRUD
		users := admin.Group("/users")
		{
			// List (Search, Filter, Pagination)
			users.GET("", handler.ListUsers)
			// Daftar role dinamis
			users.GET("/roles", handler.GetRoles)
			// Detail user
			users.GET("/:id", handler.GetUserByID)
			// Tambah baru
			users.POST("", handler.CreateUser)
			// Update (Partial)
			users.PATCH("/:id", handler.UpdateUser)
			// Toggle Aktif/Nonaktif
			users.PATCH("/:id/status", handler.ToggleStatus)
			// Hapus
			users.DELETE("/:id", handler.DeleteUser)
		}
	}
}
