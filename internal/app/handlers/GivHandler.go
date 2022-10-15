package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GivHandler
// Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
// и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location
func GivHandler(c *gin.Context) {
	key := c.Param("id")
	if key == "" {
		http.Error(c.Writer, "Ошибка задан пустой id", http.StatusBadRequest)
		return
	}
	get, err := constans.GetLinksStorage().Get(key)
	if err != nil {
		log.Println(constans.ErrorReadStorage, c.Request.URL, err.Error())
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	if get == "" {
		http.Error(c.Writer, "Ошибка по ключу ничего не нашлось", http.StatusBadRequest)
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Header().Add("Location", get)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
