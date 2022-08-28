package main

import (
	"ServiceShortLinks/internal/app/handlers"
	"ServiceShortLinks/internal/app/server"
	"ServiceShortLinks/internal/constans"
	mem_storage "ServiceShortLinks/internal/storage/memstorage"
	"errors"
	"log"
)

func main() {

	storage, err := mem_storage.Init()
	if err != nil {
		log.Fatalln(errors.New("Ошибка иницилизации mem_storage " + err.Error()))
	}
	//иницилизирум глобальное хранилище
	constans.GlobalStorage = mem_storage.MemStorage{Connect: storage}
	router := handlers.Router()
	err = server.Server(router)
	if err != nil {
		log.Fatalln(err)
	}
}
