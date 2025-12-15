package config

import (
	"os"

	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURI string
	ImagePath   string
	StoragePath string
	Port        string
	Environment string
}

func LoadConfig() *Config {
	var err = godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	return &Config{
		DatabaseURI: os.Getenv("DATABASE_URI"),
		ImagePath:   os.Getenv("IMAGE_PATH"),
		Port:        os.Getenv("HTTP_PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		StoragePath: os.Getenv("STORAGE_PATH"),
	}
}
