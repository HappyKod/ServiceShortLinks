// Package utils вспомогательные функции.
package utils

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
)

var re = regexp.MustCompile(`(\b(https?)://)?[-A-Za-z0-9+&@#/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#/%=~_|]+`)

// ValidatorURL валидирует ссылку.
func ValidatorURL(rawText string) bool {
	return re.MatchString(rawText)
}

// GeneratorStringUUID создает уникальный uuid.
func GeneratorStringUUID() string {
	return uuid.New().String()
}

// GenerateURL создает ссылку-ключ.
func GenerateURL(key string) (string, error) {
	return url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
}

// ResolveIP парсим ip.
func ResolveIP(r *http.Request) (net.IP, error) {
	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}

	if ip == nil {
		addr := r.RemoteAddr
		ipStr2, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ip = net.ParseIP(ipStr2)
		if ip == nil {
			return nil, fmt.Errorf("unexpected parse ip error")
		}
	}

	if ip == nil {
		return nil, fmt.Errorf("failed parse ip from http header")
	}

	return ip, nil

}
