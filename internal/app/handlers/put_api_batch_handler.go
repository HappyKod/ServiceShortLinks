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
	"time"
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

	var links []models.Link
	for _, v := range bodyRequest {
		if !utils.ValidatorURL(v.URL) {
			http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
			return
		}
		links = append(links, models.Link{
			FullURL:  v.URL,
			ShortKey: utils.GeneratorStringUUID(),
			UserID:   userID,
			Created:  time.Now(),
		})
	}
	if err = linksStorage.ManyPutShortLink(links); err != nil {
		log.Println(constans.ErrorWriteStorage, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
		return
	}
	var body []map[string]string
	for _, link := range links {
		for _, v := range bodyRequest {
			if v.URL == link.FullURL {
				shortURL, err := utils.GenerateURL(link.ShortKey)
				if err != nil {
					log.Println(constans.ErrorGenerateURL, link.ShortKey, err)
					http.Error(c.Writer, constans.ErrorGenerateURL, http.StatusInternalServerError)
					return
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
		log.Println(constans.ErrorWriteBody, body, err)
		http.Error(c.Writer, constans.ErrorWriteBody, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
	c.Header("content-type", "application/json")
	_, err = c.Writer.Write(bytes)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, string(bytes), err.Error())
	}
}
