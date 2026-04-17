package main

import (
	"fmt"
	"log"

	"github.com/ahmadzakyarifin/school-payment-system/config"
	"github.com/ahmadzakyarifin/school-payment-system/internal/infrastructure"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/password"
)

func main() {
	cfg, err := config.ConfigDB()
	if err != nil {
		log.Fatal(err)
	}

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username := "admin"
	plainPassword := "admin123"

	hashed, _ := password.Hash(plainPassword)

	query := `
		INSERT INTO users (name, username, password_hash, role, is_active)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash)
	`

	_, err = db.Exec(query, "Administrator", username, hashed, "admin", true)
	if err != nil {
		log.Fatalf("Gagal membuat admin: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("Admin Account Created Successfully!")
	fmt.Println("Username: ", username)
	fmt.Println("Password: ", plainPassword)
	fmt.Println("========================================")
}
