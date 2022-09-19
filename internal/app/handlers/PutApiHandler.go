package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage"
	"HappyKod/ServiceShortLinks/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path"
)

// PutAPIHandler принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<shorten_url>"}.
func PutAPIHandler(c *gin.Context) {
	connect := constans.GlobalContainer.Get("links-storage").(storage.Storages)
	bytesStructURL, err := io.ReadAll(c.Request.Body)
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
	var bodyRequest struct {
		URL string `xml:"url"`
	}
	err = json.Unmarshal(bytesStructURL, &bodyRequest)
	if err != nil {
		if err != nil {
			log.Println("Ошибка преобразования тела запроса ", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
			return
		}
	}
	if !utils.ValidatorURL(bodyRequest.URL) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}

	key, err := connect.CreateUniqKey()
	if err != nil {
		log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		return
	}
	if err = connect.Put(key, bodyRequest.URL); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	url := constans.GlobalContainer.Get("server-config").(models.Config).BaseURL
	fmt.Println(url, "-----")
	bodyResponse := struct {
		Result string `xml:"result"`
	}{
		"http://" + path.Join(url, key),
	}

	bytes, err := json.Marshal(bodyResponse)
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytes), key, err.Error())
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	_, err = c.Writer.Write(bytes)
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(bytes), key, err.Error())
		return
	}

}
