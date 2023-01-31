package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/models"

	"github.com/caarlos0/env/v6"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func version() {
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)
}
func main() {
	var cfg models.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err, "Ошибка считывания конфига")
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "адрес запуска HTTP-сервера")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePATH, "f", cfg.FileStoragePATH, "путь до файла с сокращёнными URL")
	flag.StringVar(&cfg.DataBaseURL, "d", cfg.DataBaseURL, "строка с адресом подключения к БД")
	flag.StringVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "включения HTTPS в веб-сервере")
	flag.StringVar(&cfg.FileCONFIG, "c", cfg.FileCONFIG, "возможность конфигурации приложения с помощью файла в формате JSON.")
	flag.Parse()

	if cfg.FileCONFIG != "" {
		open, err := os.Open(cfg.FileCONFIG)
		if err != nil {
			log.Fatal("ошибка считывания текстового конфига", err)
		}
		all, err := io.ReadAll(open)
		if err != nil {
			log.Fatal("ошибка считывания текстового конфига", err)
		}
		err = json.Unmarshal(all, &cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = open.Close()
		if err != nil {
			log.Println("ошибка закрытия файла", cfg.FileCONFIG)
		}
	}
	version()
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	router := handlers.Router()
	server.NewServer(router)
}
