package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/school-payment-system/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(cfg *config.Config) (*sqlx.DB,error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",cfg.DBUser,cfg.DBPass,cfg.DBHost,cfg.DBPort,cfg.DBName)

	db, err := sqlx.Connect("mysql",dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	// jumlah connect yang dapat bekerja bersamaan
	db.SetMaxOpenConns(25)
	// jumlah connect stanby
	db.SetMaxIdleConns(10)
	// max life connect stanby sampai di reycle
	db.SetConnMaxIdleTime(10 * time.Minute)
	// max life connect bekerja  sampai di reycle
	db.SetConnMaxLifetime(30 * time.Minute)	

	// PING: Verifikasi apakah database benar-benar bisa dijangkau
	// Pakai Context dengan timeout agar aplikasi tidak "hang" kelamaan kalau DB mati
	ctx,cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database tidak merespon (ping gagal): %w", err)
	}

	return db, nil
}