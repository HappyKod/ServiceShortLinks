// Package models конфигурация приложения.
package models

// Config стартовая конфигурация приложения.
type Config struct {
	Address         string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePATH string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"https://github.com/HappyKod/ServiceShortLinks"`
	DataBaseURL     string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     string `env:"ENABLE_HTTPS" json:"enable_https"`
	FileCONFIG      string `env:"CONFIG"`
}
