package handlers

import (
	"ServiceShortLinks/internal/constans"
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
	log.Println("Получен запрос на извелечение url", c.Request.URL, key)
	get, err := constans.GlobalStorage.Get(key)
	if err != nil {
		log.Println("Ошибка получение данных из хранилища ", c.Request.URL, err.Error())
		http.Error(c.Writer, "Ошибка получение данных из хранилища ", http.StatusInternalServerError)
		return
	}
	if get == "" {
		http.Error(c.Writer, "Ошибка по ключу ничего не нашлось", http.StatusBadRequest)
		return
	}
	c.Writer.Header().Set("Location", get)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	log.Println("Данные получены по ", c.Request.URL, get)
}
