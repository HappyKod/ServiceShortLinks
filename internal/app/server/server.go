package server

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// NewServer создания сервера с настройками
func NewServer(r *gin.Engine) {
	server := http.Server{
		Handler: r,
		Addr:    constans.Address,
	}
	log.Fatalln(server.ListenAndServe())
}
