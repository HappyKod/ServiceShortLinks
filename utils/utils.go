package utils

import (
	"github.com/google/uuid"
	"regexp"
)

//ValidatorUrl валидирует ссылку
func ValidatorUrl(rawText string) bool {
	var re = regexp.MustCompile(`(\b(https?):\/\/)?[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+\.[-A-Za-z0-9+&@#\/%=~_|]+`)
	return re.Match([]byte(rawText))
}

//GeneratorStringUuid создает уникальный uuid
func GeneratorStringUuid() string {
	return uuid.New().String()
}
