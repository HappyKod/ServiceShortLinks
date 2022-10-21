package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// Router указание маршрутов хендрлеров
func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.GzipWriter())
	r.Use(middleware.GzipReader())
	r.Use(middleware.WorkCooke())
	r.GET("/:id", GivHandler)
	r.POST("/", PutHandler)
	r.GET("/ping", PingDataBaseHandler)
	groupAPI := r.Group("/api")
	{
		groupAPI.GET("/user/urls", GivUsersLinksHandler)
		groupAPI.POST("/shorten", PutAPIHandler)
		groupAPI.POST("/shorten/batch", PutAPIBatchHandler)

	}
	return r
}
