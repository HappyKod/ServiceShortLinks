package constans

import "errors"

const (
	ErrorReadBody       = "ошибка при обработки тела запроса"
	ErrorWriteBody      = "ошибка при записи тела запроса"
	ErrorCloseBody      = "ошибка при закрытии тела запроса"
	ErrorReadStorage    = "ошибка при чтении данных из хранилища"
	ErrorGetKeyStorage  = "ошибка при получении значение ключа по полной ссылке"
	ErrorWriteStorage   = "ошибка при записи ссылки в хранилище"
	ErrorConnectStorage = "ошибка при попытки создать подключение"
	ErrorInvalidUrl     = "ошибка ссылка не валидна"
	ErrorGenerateUrl    = "ошибка генерации ссылки"
)

var (
	ErrorNotFindFullUrl  = errors.New("ошибка полная ссылка не найдена")
	ErrorNoUNIQUEFullUrl = errors.New("ошибка полная ссылка не уникальна")
)
