package handlers

import (
	"ServiceShortLinks/internal/constans"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router указание маршртуов сревера
func Router() *http.Server {
	r := gin.New()
	r.GET("/:id", func(context *gin.Context) { GivHandler(context) })
	r.POST("/", func(context *gin.Context) { PutHandler(context) })
	server := http.Server{
		Handler: r,
		Addr:    constans.Adres,
	}
	return &server
}
