package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GzipReader Обработчик gzip сжатия
func GzipReader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get(`Content-Encoding`) != `gzip` {
			c.Next()
			return
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Ошибка обработки сжатого тела запроса gzip", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка обработки сжатого тела запроса gzip", http.StatusInternalServerError)
			return
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Println("Ошибка не получилось закрыть Body", err)
			}
		}(c.Request.Body)
		reader, err := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
		if err != nil {
			log.Println("Ошибка обработки сжатого тела запроса gzip", c.Request.URL, err.Error())
			http.Error(c.Writer, "Ошибка обработки сжатого тела запроса gzip", http.StatusInternalServerError)
			return
		}
		c.Request.Body = reader
		c.Next()
	}
}
