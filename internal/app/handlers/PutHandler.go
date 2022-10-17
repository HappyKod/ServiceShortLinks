package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/url"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	usersStorage := constans.GetUsersStorage()
	userID := c.Param(constans.CookeUserIDName)
	bytesURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	defer func() {
		err = c.Request.Body.Close()
		if err != nil {
			log.Println(constans.ErrorCloseBody, err)
			return
		}
	}()
	if !utils.ValidatorURL(string(bytesURL)) {
		http.Error(c.Writer, "ошибка ссылка не валидна", http.StatusBadRequest)
		return
	}
	key, err := linksStorage.CreateUniqKey()
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	if err = linksStorage.PutShortLink(key, string(bytesURL)); err != nil {
		if err.(*pq.Error).Code == pgerrcode.UniqueViolation {
			getKey, err := linksStorage.GetKey(string(bytesURL))
			if err != nil {
				log.Println(constans.ErrorReadStorage, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusConflict)
			uri, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, getKey)
			if err != nil {
				log.Println("ошибка генерации ссылки", c.Request.URL, getKey, err)
				http.Error(c.Writer, "ошибка генерации ссылки", http.StatusInternalServerError)
			}
			_, err = c.Writer.WriteString(uri)
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
	c.Writer.WriteHeader(http.StatusCreated)
	body, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
	if err != nil {
		log.Println("ошибка генерации ссылки", c.Request.URL, string(bytesURL), key, err)
	}
	_, err = c.Writer.WriteString(body)
	if err != nil {
		log.Println("ошибка генерации Body ", c.Request.URL, string(bytesURL), key, err)
	}
}
