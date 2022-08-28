package utils

import (
	"github.com/google/uuid"
	"regexp"
)

// ValidatorURL валидирует ссылку
func ValidatorURL(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

// GeneratorStringUUID создает уникальный uuid
func GeneratorStringUUID() string {
	return uuid.New().String()
}
