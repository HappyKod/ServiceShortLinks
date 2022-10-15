package middleware

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

const errorReadCooke = "ошибка считывания cooke"
const errorValidCooke = "ошибка валидации cooke"

// WorkCooke Обработчик cooke
func WorkCooke() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(constans.CookeSessionName)
		if err != nil {
			if !strings.EqualFold(err.Error(), http.ErrNoCookie.Error()) {
				log.Println(errorReadCooke, err)
				http.Error(c.Writer, errorReadCooke, http.StatusBadGateway)
				return
			}
			generateCookie(c)
			c.Next()
			return
		}
		valid, err := validCookie(c, cookie)
		if err != nil {
			log.Println(errorValidCooke, cookie, err)
		}
		if !valid {
			generateCookie(c)
		}
		c.Next()
	}
}

// generateCookie генерируем новую cooke
func generateCookie(c *gin.Context) {
	h := hmac.New(sha256.New, constans.GlobalContainer.Get("secret-key").([]byte))
	userID := []byte(utils.GeneratorStringUUID()[:16])
	h.Write(userID)
	dst := h.Sum(nil)
	var cooke []byte
	cooke = append(cooke, userID[:]...)
	cooke = append(cooke, dst...)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(c.Writer, &http.Cookie{Name: constans.CookeSessionName,
		Value:   fmt.Sprintf("%x", cooke),
		Expires: expiration, Path: "/"})
	c.AddParam(constans.CookeUserIDName, string(userID))
}

// validCookie проверка cooke
func validCookie(c *gin.Context, cooke string) (bool, error) {
	data, err := hex.DecodeString(cooke)
	if err != nil {
		return false, err
	}
	if len(data[:]) < 16 {
		return false, errors.New("длина cooke не соответствует требованиям")
	}
	h := hmac.New(sha256.New, constans.GlobalContainer.Get("secret-key").([]byte))
	h.Write(data[:16])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[16:]) {
		c.AddParam(constans.CookeUserIDName, string(data[:16]))
		return true, nil
	}
	return false, nil
}
