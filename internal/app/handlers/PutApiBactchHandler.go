package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
)

// PutAPIBatchHandler принимающий в теле запроса множество URL для сокращения в формате:
//
//				[
//			   {
//			       "correlation_id": "<строковый идентификатор>",
//			       "original_url": "<URL для сокращения>"
//			   },
//			   ...
//		 ]
//		 В качестве ответа хендлер должен возвращать данные в формате:
//	 [
//	   {
//	       "correlation_id": "<строковый идентификатор из объекта запроса>",
//	       "short_url": "<результирующий сокращённый URL>"
//	   },
//	   ...
//	 ]
func PutAPIBatchHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	usersStorage := constans.GetUsersStorage()
	userID := c.Param(constans.CookeUserIDName)
	bytesStructURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = c.Request.Body.Close(); err != nil {
			log.Println(constans.ErrorCloseBody, err)
		}
	}()
	var bodyRequest []struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}
	err = json.Unmarshal(bytesStructURL, &bodyRequest)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}

	var urls []string
	for _, v := range bodyRequest {
		if !utils.ValidatorURL(v.URL) {
			http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
			return
		}
		urls = append(urls, v.URL)
	}

	result, err := linksStorage.ManyPut(urls)
	if err != nil {
		log.Println(constans.ErrorWriteStorage, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
		return
	}
	go func() {
		for k, _ := range result {
			if err = usersStorage.Put(userID, k); err != nil {
				log.Println(constans.ErrorWriteStorage, c.Request.URL, err.Error())
			}
		}
	}()
	var body []map[string]string
	for key, uri := range result {
		for _, v := range bodyRequest {
			if v.URL == uri {
				shortURL, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
				if err != nil {
					log.Println("Ошибка генерации ссылки", c.Request.URL, key, err.Error())
					http.Error(c.Writer, "Ошибка генерации ссылки", http.StatusInternalServerError)

				}
				body = append(body, map[string]string{
					"correlation_id": v.ID,
					"short_url":      shortURL,
				})
				break
			}
		}
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, string(bytes), err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	_, err = c.Writer.Write(bytes)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, string(bytes), err.Error())
	}
}
