// Package server запуск сервера.
package server

import (
	"log"
	"net/http"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// NewServer создания сервера с настройками.
func NewServer(r *gin.Engine) {
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	if cfg.EnableHTTPS != "" {
		// конструируем менеджер TLS-сертификатов
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}
		// конструируем сервер с поддержкой TLS
		server := http.Server{
			Addr:    ":443",
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
		log.Fatalln(server.ListenAndServeTLS("", ""))
	} else {
		server := http.Server{
			Handler: r,
			Addr:    cfg.Address,
		}
		log.Fatalln(server.ListenAndServe())
	}
}
