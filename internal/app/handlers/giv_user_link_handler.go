/*
Package handlers работа GivUsersLinksHandler возвращает пользователю
все когда-либо сокращённые им URL.
*/
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/utils"

	"github.com/gin-gonic/gin"
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
	userID := c.Param(constans.CookeUserIDName)
	var doneLinks []map[string]string
	links, err := linksStorage.GetShortLinkUser(userID)
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	for _, link := range links {
		shortLink, err := utils.GenerateURL(link.ShortKey)
		if err != nil {
			log.Println(constans.ErrorGenerateURL, link.ShortKey, err)
			http.Error(c.Writer, constans.ErrorGenerateURL, http.StatusInternalServerError)
			return
		}
		doneLinks = append(doneLinks, map[string]string{"short_url": shortLink, "original_url": link.FullURL})
	}
	if len(doneLinks) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	body, err := json.Marshal(doneLinks)
	if err != nil {
		log.Println(constans.ErrorWriteBody, c.Request.URL, doneLinks, err)
		http.Error(c.Writer, constans.ErrorWriteBody, http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("content-type", "application/json")
	_, err = c.Writer.Write(body)
	if err != nil {
		log.Println(constans.ErrorWriteBody, string(body), err)
	}
}
