package server

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewServer создания сервера с настройками
func NewServer(r *gin.Engine) *http.Server {
	server := http.Server{
		Handler: r,
		Addr:    constans.Address,
	}
	return &server
}

// Server запуск сервера
func Server(server *http.Server) error {
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
