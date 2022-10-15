package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
)

// GivUsersLinksHandler хендлер который сможет вернуть пользователю
//
//	все когда-либо сокращённые им URL в формате
//	[
//	  {
//	      "short_url": "http://...",
//	      "original_url": "http://..."
//	  },
//	  ...
//	]
func GivUsersLinksHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	usersStorage := constans.GetUsersStorage()
	userID := c.Param(constans.CookeUserIDName)
	var doneLinks []map[string]string
	links, err := usersStorage.Get(userID)
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	for _, key := range links {
		fullLink, err := linksStorage.Get(key)
		if err != nil {
			log.Println(constans.ErrorReadStorage, c.Request.URL, err)
			http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
			return
		}
		shortLink, err := url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
		if err != nil {
			log.Println("Ошибка генерации ссылки", c.Request.URL, key, err)
		}
		doneLinks = append(doneLinks, map[string]string{"short_url": shortLink, "original_url": fullLink})
	}
	if len(doneLinks) == 0 {
		http.Error(c.Writer, "Сокращенных ссылок у данного пользователя не найдено", http.StatusNoContent)
		return
	}
	body, err := json.Marshal(doneLinks)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, doneLinks, err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	_, err = c.Writer.Write(body)
	if err != nil {
		log.Println("Ошибка генерации Body ", c.Request.URL, string(body), err)
	}
}
