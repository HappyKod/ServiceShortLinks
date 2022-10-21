package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage/pglinkssotorage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
