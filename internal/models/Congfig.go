// Package models конфигурация приложения
package models

// Config стартовая конфигурация приложения
type Config struct {
	Address         string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePATH string `env:"FILE_STORAGE_PATH"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"https://github.com/HappyKod/ServiceShortLinks"`
	DataBaseURL     string `env:"DATABASE_DSN"`
}
