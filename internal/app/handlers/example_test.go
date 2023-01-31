package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"

	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
)

func ExamplePutHandler() {
	// Поднимаем Конфигурацию.
	cfg := models.Config{}
	// Поднимаем Контейнер.
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	router := Router()
	w := httptest.NewRecorder()
	// Создаем запрос.
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://github.com/HappyKod/ServiceShortLinks")))
	// Совершаем запрос.
	router.ServeHTTP(w, req)
}

func ExampleGivHandler() {
	//Поднимаем Конфигурацию.
	cfg := models.Config{}
	key := "test1"
	//Поднимаем Контейнер.
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	link := models.Link{ShortKey: key, FullURL: "https://github.com/HappyKod/ServiceShortLinks"}
	err = constans.GetLinksStorage().PutShortLink(key, link)
	if err != nil {
		log.Fatal(err)
	}
	router := Router()
	w := httptest.NewRecorder()
	//Создаем запрос.
	req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
	//Совершаем запрос.
	router.ServeHTTP(w, req)
}
