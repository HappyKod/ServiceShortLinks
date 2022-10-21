package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие byte, поэтому пишем в него
	return w.writer.Write(b)
}

func (w gzipWriter) WriteString(s string) (int, error) {
	// w.WriteString будет отвечать за gzip-сжатие string, поэтому пишем в него
	w.Header().Del("Content-Length")
	return w.writer.Write([]byte(s))
}

// GzipWriter Обработчик gzip сжатия
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
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		c.Writer = &gzipWriter{c.Writer, gz}
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Next()
	}
}
