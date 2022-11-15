package constans

import "errors"

const (
	ErrorReadBody       = "ошибка при обработки тела запроса"
	ErrorWriteBody      = "ошибка при записи тела запроса"
	ErrorCloseBody      = "ошибка при закрытии тела запроса"
	ErrorReadStorage    = "ошибка при чтении данных из хранилища"
	ErrorGetKeyStorage  = "ошибка при получении значение ключа по полной ссылке"
	ErrorWriteStorage   = "ошибка при записи ссылки в хранилище"
	ErrorUpdateStorage  = "ошибка при обновлении данных в хранилище"
	ErrorConnectStorage = "ошибка при попытки создать подключение"
	ErrorInvalidURL     = "ошибка ссылка не валидна"
	ErrorGenerateURL    = "ошибка генерации ссылки"
)

var (
	ErrorNotFindFullURL  = errors.New("ошибка полная ссылка не найдена")
	ErrorNoUNIQUEFullURL = errors.New("ошибка полная ссылка не уникальна")
)
