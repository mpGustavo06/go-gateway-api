package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/repository"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/service"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/web/server"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf (
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	port := getEnv("PORT", "8080")
	server := server.NewServer(accountService, port)
	
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server", err)
	}
}
