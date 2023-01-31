/*
Package handlers работа PingDataBaseHandler
проверяет соединение с базой данных.
*/
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HappyKod/ServiceShortLinks/internal/constans"
	"github.com/HappyKod/ServiceShortLinks/internal/models"
	"github.com/HappyKod/ServiceShortLinks/internal/storage/linksstorage/pglinkssotorage"
)

// PingDataBaseHandler проверяет соединение с базой данных.
func PingDataBaseHandler(c *gin.Context) {
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	linkStorage, err := pglinkssotorage.New(cfg.DataBaseURL)
	if err != nil {
		log.Println(constans.ErrorConnectStorage, err)
		http.Error(c.Writer, constans.ErrorConnectStorage, http.StatusInternalServerError)
		return
	}
	if err = linkStorage.Ping(); err != nil {
		log.Println(constans.ErrorConnectStorage, err)
		http.Error(c.Writer, constans.ErrorConnectStorage, http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
