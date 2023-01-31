// Package server запуск сервера.
package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HappyKod/ServiceShortLinks/internal/constans"
	"github.com/HappyKod/ServiceShortLinks/internal/models"
)

// NewServer создания сервера с настройками.
func NewServer(r *gin.Engine) {
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	server := http.Server{
		Handler: r,
		Addr:    cfg.Address,
	}
	log.Fatalln(server.ListenAndServe())
}
