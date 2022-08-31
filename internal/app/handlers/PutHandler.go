package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/storage"
	"HappyKod/ServiceShortLinks/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	storage := constans.GlobalContainer.Get("links-storage").(storage.Storages)
	bytesURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Ошибка обработки тела запроса ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Ошибка закрытия тела запроса ", err)
			http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
			return
		}
	}(c.Request.Body)
	if !utils.ValidatorURL(string(bytesURL)) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}
	var key string
	//Подбираем уникальный ключ
	for {
		key = utils.GeneratorStringUUID()
		get, err := storage.Get(key)
		if err != nil {
			log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
			return
		}
		if get == "" {
			break
		}
	}
	if err = storage.Put(key, string(bytesURL)); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.WriteString("http://" + path.Join(constans.Address, key))
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytesURL), key, err.Error())
		http.Error(c.Writer, "Ошибка генерации Body", http.StatusInternalServerError)
	}
}
