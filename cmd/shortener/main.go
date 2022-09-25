package main

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/models"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	var cfg models.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err, "Ошибка считывания конфига")
	}
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	router := handlers.Router()
	server.NewServer(router)
}
