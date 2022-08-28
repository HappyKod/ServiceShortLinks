package handlers

import (
	"ServiceShortLinks/internal/constans"
	"ServiceShortLinks/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	bytesUrl, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Ошибка обработки тела запроса ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
		return
	}
	if !utils.ValidatorURL(string(bytesUrl)) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}
	var key string
	//Подбираем уникальный ключ
	for {
		key = utils.GeneratorStringUUID()
		get, err := constans.GlobalStorage.Get(key)
		if err != nil {
			log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		}
		if get == "" {
			break
		}
	}
	if err = constans.GlobalStorage.Put(key, string(bytesUrl)); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.WriteString(key)
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytesUrl), key, err.Error())
		http.Error(c.Writer, "Ошибка генерации Body", http.StatusInternalServerError)
	}
	return
}
