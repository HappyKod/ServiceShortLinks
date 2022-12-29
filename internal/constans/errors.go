package constans

import "errors"

// Набор ошибок строкой
const (
	//ErrorReadBody ошибка при обработки тела запроса
	ErrorReadBody = "ошибка при обработки тела запроса"
	//ErrorWriteBody ошибка при записи тела запроса
	ErrorWriteBody = "ошибка при записи тела запроса"
	//ErrorCloseBody ошибка при закрытии тела запроса
	ErrorCloseBody = "ошибка при закрытии тела запроса"
	//ErrorReadStorage ошибка при чтении данных из хранилища
	ErrorReadStorage = "ошибка при чтении данных из хранилища"
	//ErrorGetKeyStorage ошибка при получении значение ключа по полной ссылке
	ErrorGetKeyStorage = "ошибка при получении значение ключа по полной ссылке"
	//ErrorWriteStorage ошибка при записи ссылки в хранилище
	ErrorWriteStorage = "ошибка при записи ссылки в хранилище"
	//ErrorUpdateStorage ошибка при обновлении данных в хранилище
	ErrorUpdateStorage = "ошибка при обновлении данных в хранилище"
	//ErrorConnectStorage ошибка при попытки создать подключение
	ErrorConnectStorage = "ошибка при попытки создать подключение"
	//ErrorInvalidURL ошибка ссылка не валидна
	ErrorInvalidURL = "ошибка ссылка не валидна"
	//ErrorGenerateURL ошибка генерации ссылки
	ErrorGenerateURL = "ошибка генерации ссылки"
)

// Набор ошибок
var (
	// ErrorNotFindFullURL ошибка полная ссылка не найдена
	ErrorNotFindFullURL = errors.New("ошибка полная ссылка не найдена")
	//ErrorNoUNIQUEFullURL ошибка полная ссылка не уникальна
	ErrorNoUNIQUEFullURL = errors.New("ошибка полная ссылка не уникальна")
)
