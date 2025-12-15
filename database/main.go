package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run cmd/migrate/main.go [up|down]")
	}

	command := os.Args[1]

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	// database connection string
	dbURL := os.Getenv("DATABASE_URI")

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		logger.Fatal("Error open database connection")
	}

	// set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	migrationsDir := "database/migration"

	switch command {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migration up success")

	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migration down success")

	default:
		log.Fatalf("unknown command: %s (use up or down)", command)
	}
}
