// Package handlers работа PutAPIHandler возвращает в теле запроса в формате JSON исходный URL.
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/HappyKod/ServiceShortLinks/internal/constans"
	"github.com/HappyKod/ServiceShortLinks/internal/models"
	"github.com/HappyKod/ServiceShortLinks/utils"
)

// PutAPIHandler принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<shorten_url>"}.
func PutAPIHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	userID := c.Param(constans.CookeUserIDName)
	bytesStructURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = c.Request.Body.Close(); err != nil {
			log.Println(constans.ErrorCloseBody, err)
		}
	}()
	var bodyRequest struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(bytesStructURL, &bodyRequest)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	if !utils.ValidatorURL(bodyRequest.URL) {
		http.Error(c.Writer, "Ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}
	link := models.Link{
		ShortKey: utils.GeneratorStringUUID(),
		FullURL:  bodyRequest.URL,
		UserID:   userID,
		Created:  time.Now(),
	}
	var key string
	if err = linksStorage.PutShortLink(link.ShortKey, link); err != nil {
		if errPG, ok := err.(*pq.Error); ok {
			if errPG.Code != pgerrcode.UniqueViolation {
				log.Println(constans.ErrorWriteStorage, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
				return
			}
		} else {
			if !errors.Is(constans.ErrorNoUNIQUEFullURL, err) {
				log.Println(constans.ErrorWriteStorage, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
				return
			}
		}
		key, err = linksStorage.GetKey(bodyRequest.URL)
		if err != nil {
			log.Println(constans.ErrorGetKeyStorage, bodyRequest.URL, err)
			http.Error(c.Writer, constans.ErrorGetKeyStorage, http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusConflict)
	} else {
		key = link.ShortKey
		c.Status(http.StatusCreated)
	}
	body, err := utils.GenerateURL(key)
	if err != nil {
		log.Println(constans.ErrorGenerateURL, key, err)
		http.Error(c.Writer, constans.ErrorGenerateURL, http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(map[string]string{"result": body})
	if err != nil {
		log.Println(constans.ErrorWriteBody, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorWriteBody, http.StatusInternalServerError)
		return
	}
	c.Header("content-type", "application/json")
	_, err = c.Writer.Write(bytes)
	if err != nil {
		log.Println(constans.ErrorWriteBody, c.Request.URL, err)
	}
}
