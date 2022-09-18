package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
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

// PutApiHandler принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<shorten_url>"}.
func PutApiHandler(c *gin.Context) {
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
		Url string `xml:"url"`
	}
	err = json.Unmarshal(bytesStructURL, &bodyRequest)
	if err != nil {
		if err != nil {
			log.Println("Ошибка преобразования тела запроса ", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка обработки тела запроса", http.StatusInternalServerError)
			return
		}
	}
	if !utils.ValidatorURL(bodyRequest.Url) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}

	key, err := connect.CreateUniqKey()
	if err != nil {
		log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		return
	}
	if err = connect.Put(key, bodyRequest.Url); err != nil {
		log.Println("Ошибка записи данных в хранилище ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка записи данных в хранилище", http.StatusInternalServerError)
		return
	}
	get, _ := connect.Get(key)
	fmt.Println(get)

	bodyResponse := struct {
		Result string `xml:"result"`
	}{
		"http://" + path.Join(constans.Address, key),
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
