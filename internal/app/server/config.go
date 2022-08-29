package server

import (
	"ServiceShortLinks/internal/constans"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Config(r *gin.Engine) *http.Server {
	server := http.Server{
		Handler: r,
		Addr:    constans.Address,
	}
	return &server
}
