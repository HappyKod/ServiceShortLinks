package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage"
	"HappyKod/ServiceShortLinks/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	connect := constans.GlobalContainer.Get("links-storage").(storage.Storages)
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
	key, err := connect.CreateUniqKey()
	if err != nil {
		log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		return
	}
	if err = connect.Put(key, string(bytesURL)); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	body, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
	if err != nil {
		log.Println("Ошибка генерации ссылки", c.Request.URL, string(bytesURL), key, err.Error())
	}
	_, err = c.Writer.WriteString(body)
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytesURL), key, err.Error())
	}
}
