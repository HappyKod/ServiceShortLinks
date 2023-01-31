// Package handlers работа DelUsersLinksHandler удаление ссылок пользователя.
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"HappyKod/ServiceShortLinks/internal/constans"

	"github.com/gin-gonic/gin"
)

// DelUsersLinksHandler принимает список идентификаторов
//
//		сокращённых URL для удаления в формате:
//		[ "a", "b", "c", "d", ...]
//	 В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.
//	 Фактический результат удаления может происходить позже — каким-либо образом оповещать
//	 пользователя об успешности или неуспешности не нужно.
//	 Успешно удалить URL может пользователь, его создавший.
//	 При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.
func DelUsersLinksHandler(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	userID := c.Param(constans.CookeUserIDName)
	bytesBody, err := io.ReadAll(c.Request.Body)
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
	var links []string
	err = json.Unmarshal(bytesBody, &links)
	if err != nil {
		log.Println(constans.ErrorReadBody, c.Request.URL, err)
		http.Error(c.Writer, constans.ErrorReadBody, http.StatusInternalServerError)
		return
	}
	go func() {
		err = linksStorage.DeleteShortLinkUser(userID, links)
		if err != nil {
			log.Println(constans.ErrorUpdateStorage, c.Request.URL, err)
		}
	}()
	c.String(http.StatusAccepted, "Ваш запрос на удаление успешно принят")
}
