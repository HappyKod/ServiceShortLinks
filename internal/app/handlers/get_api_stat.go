package handlers

import (
	"log"
	"net"
	"net/http"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/utils"

	"github.com/gin-gonic/gin"
)

// GetStat получаем информацию о кол-во пользователей и ссылок
// возвращающий в ответ объект:
// {
// "urls": <int>, // количество сокращённых URL в сервисе
// "users": <int> // количество пользователей в сервисе
// }
func GetStat(c *gin.Context) {
	linksStorage := constans.GetLinksStorage()
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	ipNet := constans.GlobalContainer.Get("ip-net").(net.IPNet)
	ip, err := utils.ResolveIP(c.Request)
	if err != nil {
		log.Println(constans.ErrorReadStorage, err)
		http.Error(c.Writer, "", http.StatusInternalServerError)
		return
	}
	if cfg.TrustedSubnet != "" {
		if !ipNet.Contains(ip) {
			c.String(http.StatusForbidden, "")
			return
		}
	}
	stat := make(map[string]int)
	countUser, countLink, err := linksStorage.Stat()
	if err != nil {
		log.Println(constans.ErrorReadStorage, err)
		http.Error(c.Writer, constans.ErrorReadStorage, http.StatusInternalServerError)
		return
	}
	stat["users"] = countUser
	stat["urls"] = countLink
	c.JSON(http.StatusOK, stat)
}
