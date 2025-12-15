package main

import (
	"context"
	"sync"

	"github.com/dona-dllollin/belajar-clean-arch/internal/config"
	"github.com/dona-dllollin/belajar-clean-arch/internal/delivery/http"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/validation"
	"github.com/jackc/pgx/v5"
)

var wg sync.WaitGroup

func main() {

	// load Config
	cfg := config.LoadConfig()

	// initialize logger
	logger.Initialize(cfg.Environment)

	// initialize validator
	validator := validation.New()

	// database connection
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	logger.Info("Successfully connected to the database")

	httpServer := http.NewServer(validator, conn)

	wg.Add(1)

	// Run HTTP server
	go func() {
		defer wg.Done()
		if err := httpServer.Run(); err != nil {
			logger.Fatal("Running HTTP server error:", err)
		}
	}()

	wg.Wait()
}
