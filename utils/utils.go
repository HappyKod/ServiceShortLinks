package utils

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"github.com/google/uuid"
	"net/url"
	"regexp"
)

// ValidatorURL валидирует ссылку
func ValidatorURL(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?)://)?[-A-Za-z0-9+&@#/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#/%=~_|]+`)
	return re.MatchString(rawText)
}

// GeneratorStringUUID создает уникальный uuid
func GeneratorStringUUID() string {
	return uuid.New().String()
}

// GenerateURL создает ссылку-ключ
func GenerateURL(key string) (string, error) {
	return url.JoinPath(constans.GlobalContainer.Get("server-config").(models.Config).BaseURL, key)
}
