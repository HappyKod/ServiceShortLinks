package main

import (
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/memstorage"
	"github.com/caarlos0/env/v6"
	"github.com/sarulabs/di"
	"log"
)

func init() {
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			storage, err := memstorage.New()
			if err != nil {
				log.Fatalln("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: storage}, nil
		}})
	if err != nil {
		log.Fatalln("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		log.Fatalln("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()
}

func main() {
	router := handlers.Router()
	server.NewServer(router)
}
