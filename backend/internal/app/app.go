package app

import (
	"fmt"
	"net/http"

	"github.com/ahmadzakyarifin/school-payment-system/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Config *config.Config
	DB     *sqlx.DB
	Router *gin.Engine
}

func New(cfg *config.Config, db *sqlx.DB) *App {
	//  Atur mode Gin (Release atau Debug) secara dinamis
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Keamanan: Batasi proxy yang dipercaya
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	// Custom 404 — tampilkan method + path agar developer tahu URL mana yang salah
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"success": false,
			"message": fmt.Sprintf("Route [%s] %s tidak ditemukan. Pastikan URL dan method sudah benar.", c.Request.Method, c.Request.URL.Path),
		})
	})

	// Custom 405 — method ada tapi salah (misal GET ke endpoint POST)
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":    405,
			"success": false,
			"message": fmt.Sprintf("Method [%s] tidak diizinkan untuk %s.", c.Request.Method, c.Request.URL.Path),
		})
	})

	return &App{
		Config: cfg,
		DB:     db,
		Router: r,
	}
}

func (a *App) Run() error {
	return a.Router.Run(":" + a.Config.Port)
}
