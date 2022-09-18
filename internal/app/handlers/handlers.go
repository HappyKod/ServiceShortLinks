package handlers

import (
	"github.com/gin-gonic/gin"
)

// Router указание маршртуов сревера
func Router() *gin.Engine {
	r := gin.New()
	r.GET("/:id", GivHandler)
	r.POST("/", PutHandler)
	r.POST("/api/shorten", PutApiHandler)
	return r
}
