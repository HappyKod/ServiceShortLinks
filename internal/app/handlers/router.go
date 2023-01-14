// Package handlers указание всех маршрутов хендлеров и подключение middleware
package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Router указание маршрутов хендлеров
func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	pprof.Register(r)
	r.Use(middleware.GzipWriter())
	r.Use(middleware.GzipReader())
	r.Use(middleware.WorkCooke())
	r.GET("/:id", GivHandler)
	r.POST("/", PutHandler)
	r.GET("/ping", PingDataBaseHandler)
	groupAPI := r.Group("/api")
	{
		groupAPI.GET("/user/urls", GivUsersLinksHandler)
		groupAPI.DELETE("/user/urls", DelUsersLinksHandler)
		groupAPI.POST("/shorten", PutAPIHandler)
		groupAPI.POST("/shorten/batch", PutAPIBatchHandler)

	}
	return r
}
