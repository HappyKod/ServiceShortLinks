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
	r.Use(middleware.WorkCooke())
	r.GET("/:id", GivHandler)
	r.POST("/", PutHandler)
	groupAPI := r.Group("/api")
	{
		groupAPI.GET("/user/urls", GivUsersLinksHandler)
		groupAPI.POST("/shorten", PutAPIHandler)
	}
	return r
}
