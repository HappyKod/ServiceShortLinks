// Package server запуск сервера.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// timeOutShutdownService таймаут остановки сервера
const timeOutShutdownService = time.Duration(5) * time.Second

// NewServer создания сервера с настройками.
func NewServer(r *gin.Engine) {
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	server := http.Server{
		Handler: r,
		Addr:    cfg.Address,
	}
	go func() {
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
			server.Addr = ":443"
			// для TLS-конфигурации используем менеджер сертификатов
			server.TLSConfig = manager.TLSConfig()
			log.Fatalln(server.ListenAndServeTLS("", ""))
		} else {
			log.Fatalln(server.ListenAndServe())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), timeOutShutdownService)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("экстренное выключение сервиса", err)
	}
	log.Println("сервис выключен")

}
