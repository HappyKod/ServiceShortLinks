package main

import (
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/storage/memstorage"
	"github.com/sarulabs/di"
	"log"
)

func main() {
	builder, _ := di.NewBuilder()
	err := builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			storage, err := memstorage.Init()
			if err != nil {
				log.Fatalln("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: storage}, nil
		}})
	if err != nil {
		log.Fatalln("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()

	//иницилизирум глобальное хранилище
	router := handlers.Router()
	server.NewServer(router)
}
