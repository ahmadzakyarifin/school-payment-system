package routes

import (
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	"github.com/ahmadzakyarifin/school-payment-system/internal/middleware"
	adminrouter "github.com/ahmadzakyarifin/school-payment-system/internal/routes/admin"
	authrouter "github.com/ahmadzakyarifin/school-payment-system/internal/routes/auth"
	parentrouter "github.com/ahmadzakyarifin/school-payment-system/internal/routes/parent"
	"github.com/gin-contrib/cors"
)

func Setup(a *app.App) {
	a.Router.Use(cors.Default())

	v1 := a.Router.Group("/api/v1")

	authrouter.Setup(a, v1)

	protected := v1.Group("")
	protected.Use(middleware.JWT(a.Config.JWTSecret))
	{
		adminrouter.Setup(a, protected) 
		parentrouter.Setup(a, protected) 
	}
}
