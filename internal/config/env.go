package env_config

import (
	"os"
)

type Config struct {
	Main     Main
	Database Database
	Auth     Auth
}
type Database struct {
	DB_host     string
	DB_port     string
	DB_user     string
	DB_password string
	DB_name     string
}

type Main struct {
	Port string
}

type Auth struct {
	JwtKey string
}

func LoadConfig() *Config {

	//err := godotenv.Load()
	// if err != nil {
	// 	slog.Error("Error loading .env file", "error", err)
	// }
	return &Config{

		Main: Main{
			Port: os.Getenv("Port"),
		},

		Database: Database{
			DB_host:     os.Getenv("DB_host"),
			DB_port:     os.Getenv("DB_port"),
			DB_user:     os.Getenv("DB_user"),
			DB_password: os.Getenv("DB_password"),
			DB_name:     os.Getenv("DB_name"),
		},
		Auth: Auth{
			JwtKey: os.Getenv("Jwt_key"),
		},
	}
}
