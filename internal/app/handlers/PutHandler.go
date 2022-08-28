package handlers

import (
	"ServiceShortLinks/internal/constans"
	"ServiceShortLinks/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	bytesURL, err := io.ReadAll(c.Request.Body)
	log.Println("Получен запрос на добавление url ", string(bytesURL))
	if err != nil {
		log.Println("Ошибка обработки тела запроса ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
		return
	}
	defer c.Request.Body.Close()
	if !utils.ValidatorURL(string(bytesURL)) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}
	var key string
	//Подбираем уникальный ключ
	for {
		key = utils.GeneratorStringUUID()
		log.Println("Сгенерирован ключ ", key, "для ", string(bytesURL))
		get, err := constans.GlobalStorage.Get(key)
		if err != nil {
			log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		}
		if get == "" {
			break
		}
	}
	if err = constans.GlobalStorage.Put(key, string(bytesURL)); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.WriteString(path.Join(constans.Adres, key))
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytesURL), key, err.Error())
		http.Error(c.Writer, "Ошибка генерации Body", http.StatusInternalServerError)
	}
}
