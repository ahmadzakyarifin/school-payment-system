package main

import (
	"log"

	"github.com/ahmadzakyarifin/school-payment-system/config"
	"github.com/ahmadzakyarifin/school-payment-system/internal/app"
	"github.com/ahmadzakyarifin/school-payment-system/internal/infrastructure"
	"github.com/ahmadzakyarifin/school-payment-system/internal/routes"
)

func main() {
	cfg, err := config.ConfigDB()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Gagal menyambung ke database: %v", err)
	}
	defer db.Close()

	a := app.New(cfg, db)

	routes.Setup(a)

	log.Printf("Server berjalan di port %s", cfg.Port)
	if err := a.Run(); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}
