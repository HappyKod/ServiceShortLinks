package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/url"
)

// PutAPIHandler принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<shorten_url>"}.
func PutAPIHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	usersStorage := constans.GetUsersStorage()
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

	key, err := linksStorage.CreateUniqKey()
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	if err = linksStorage.PutShortLink(key, bodyRequest.URL); err != nil {
		if err.(*pq.Error).Code == pgerrcode.UniqueViolation {
			getKey, err := linksStorage.GetKey(bodyRequest.URL)
			if err != nil {
				log.Println(constans.ErrorReadStorage, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusConflict)
			c.Writer.Header().Set("content-type", "application/json")
			uri, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, getKey)
			if err != nil {
				log.Println("ошибка генерации ссылки", c.Request.URL, getKey, err)
				http.Error(c.Writer, "ошибка генерации ссылки", http.StatusInternalServerError)
			}
			bytes, err := json.Marshal(map[string]string{"result": uri})
			if err != nil {
				log.Println(constans.ErrorReadBody, c.Request.URL, uri, err)
				http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
				return
			}
			_, err = c.Writer.Write(bytes)
			if err != nil {
				log.Println(constans.ErrorReadBody, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
				return
			}
			return
		}
		log.Println(constans.ErrorWriteStorage, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
		return
	}
	if err = usersStorage.Put(userID, key); err != nil {
		log.Println(constans.ErrorWriteStorage, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
		return
	}
	body, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
	if err != nil {
		log.Println("Ошибка генерации ссылки", c.Request.URL, key, err)
		http.Error(c.Writer, "Ошибка генерации ссылки", http.StatusInternalServerError)

	}
	bytes, err := json.Marshal(map[string]string{"result": body})
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, string(bytes), key, err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	_, err = c.Writer.Write(bytes)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, string(bytes), key, err)
	}
}
