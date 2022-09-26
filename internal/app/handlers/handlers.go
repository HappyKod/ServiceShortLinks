package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// Router указание маршрутов севера
func Router() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GzipWriter())
	r.Use(middleware.GzipReader())
	r.GET("/:id", GivHandler)
	r.POST("/", PutHandler)
	r.POST("/api/shorten", PutAPIHandler)
	return r
}
