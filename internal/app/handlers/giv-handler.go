// Package handlers работа GivHandler возвращает оригинальный URL.
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HappyKod/ServiceShortLinks/internal/constans"
)

// GivHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
// и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func GivHandler(c *gin.Context) {
	key := c.Param("id")
	if key == "" {
		http.Error(c.Writer, "Ошибка задан пустой id", http.StatusBadRequest)
		return
	}
	link, err := constans.GetLinksStorage().GetShortLink(key)
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	if link.FullURL == "" {
		http.Error(c.Writer, "Ошибка по ключу ничего не нашлось", http.StatusBadRequest)
		return
	}
	if link.Del {
		c.String(http.StatusGone, "данная ссылка больше не доступна")
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Header().Add("Location", link.FullURL)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
