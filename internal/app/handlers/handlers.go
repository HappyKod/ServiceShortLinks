package handlers

import (
	"github.com/gin-gonic/gin"
)

// Router указание маршртуов сревера
func Router() *gin.Engine {
	r := gin.New()
	r.GET("/:id", func(context *gin.Context) { GivHandler(context) })
	r.POST("/", func(context *gin.Context) { PutHandler(context) })
	return r
}
