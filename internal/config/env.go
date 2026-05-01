package env_config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database
}
type Database struct {
	DB_host     string
	DB_port     string
	DB_user     string
	DB_password string
	DB_name     string
}

func LoadConfig() *Config {

	err := godotenv.Load("./.env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}
	return &Config{
		Database: Database{
			DB_host:     os.Getenv("DB_host"),
			DB_port:     os.Getenv("DB_port"),
			DB_user:     os.Getenv("DB_user"),
			DB_password: os.Getenv("DB_password"),
			DB_name:     os.Getenv("DB_name"),
		},
	}
}
