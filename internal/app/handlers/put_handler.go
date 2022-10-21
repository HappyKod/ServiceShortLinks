package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"time"
)

// PutHandler принимает в теле запроса строку URL для сокращения и
// возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func PutHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
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
	fullURL := string(bytesURL)
	if !utils.ValidatorURL(fullURL) {
		http.Error(c.Writer, constans.ErrorInvalidUrl, http.StatusBadRequest)
		return
	}
	link := models.Link{
		ShortKey: utils.GeneratorStringUUID(),
		FullURL:  fullURL,
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
			if !errors.Is(constans.ErrorNoUNIQUEFullUrl, err) {
				log.Println(constans.ErrorWriteStorage, c.Request.URL, err)
				http.Error(c.Writer, constans.ErrorWriteStorage, http.StatusInternalServerError)
				return
			}
		}

		key, err = linksStorage.GetKey(fullURL)
		if err != nil {
			log.Println(constans.ErrorGetKeyStorage, fullURL, err)
			http.Error(c.Writer, constans.ErrorGetKeyStorage, http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusConflict)
	} else {
		key = link.ShortKey
		c.Status(http.StatusCreated)
	}
	uri, err := utils.GenerateURL(key)
	if err != nil {
		log.Println(constans.ErrorGenerateUrl, key, err)
		http.Error(c.Writer, constans.ErrorGenerateUrl, http.StatusInternalServerError)
		return
	}
	_, err = c.Writer.WriteString(uri)
	if err != nil {
		log.Println(constans.ErrorWriteBody, c.Request.URL, err)
	}
}
