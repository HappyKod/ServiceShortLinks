package main

import (
	"HappyKod/ServiceShortLinks/internal/app/handlers"
	"HappyKod/ServiceShortLinks/internal/app/server"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/storage/memstorage"
	"log"
)

func main() {
	storage, err := memstorage.Init()
	if err != nil {
		log.Fatalln("Ошибка иницилизации mem_storage ", err)
	}
	//иницилизирум глобальное хранилище
	constans.GlobalStorage = memstorage.MemStorage{Connect: storage}
	router := handlers.Router()
	config := server.Config(router)
	err = server.Server(config)
	if err != nil {
		log.Fatalln(err)
	}
}
