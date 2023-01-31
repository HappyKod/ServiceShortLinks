// Package middleware работа со сжатием body.
package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// Write отвечает за gzip-сжатие byte
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// WriteString отвечает за gzip-сжатие string
func (w gzipWriter) WriteString(s string) (int, error) {
	w.Header().Del("Content-Length")
	return w.writer.Write([]byte(s))
}

// GzipWriter Обработчик gzip сжатия.
func GzipWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			c.Next()
			return
		}
		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			_, err = io.WriteString(c.Writer, err.Error())
			if err != nil {
				log.Println("Ошибка не получилось записать Данные в ", err)
			}
			return
		}
		defer func() {
			if c.Writer.Size() > 0 {
				err = gz.Close()
				if err != nil {
					log.Println("Ошибка не получилось закрыть gzip.Writer", err)
				}
			}
		}()
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных.
		c.Writer = &gzipWriter{c.Writer, gz}
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Next()
	}
}
