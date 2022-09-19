package server

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// NewServer создания сервера с настройками
func NewServer(r *gin.Engine) {
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	server := http.Server{
		Handler: r,
		Addr:    cfg.Address,
	}
	log.Fatalln(server.ListenAndServe())
}
