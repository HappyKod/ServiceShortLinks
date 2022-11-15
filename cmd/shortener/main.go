package main

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/models"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	var cfg models.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err, "Ошибка считывания конфига")
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "адрес запуска HTTP-сервера")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePATH, "f", cfg.FileStoragePATH, "путь до файла с сокращёнными URL")
	flag.StringVar(&cfg.DataBaseURL, "d", cfg.DataBaseURL, "строка с адресом подключения к БД")
	flag.Parse()
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	router := handlers.Router()
	server.NewServer(router)
}
