// Package constans константы проекта.
package constans

import (
	"github.com/sarulabs/di"

	"github.com/HappyKod/ServiceShortLinks/internal/storage/linksstorage"
)

// GlobalContainer DI контейнер для общего доступа к данным.
var GlobalContainer di.Container

// Константы Сессии
const (
	CookeSessionName = "User-session" // Наименование ключа сессии.
	CookeUserIDName  = "UserID"       // Наименование ключа UserID.
	CookeUserIDLen   = 16             // Срез ключа сессии.
)

// GetLinksStorage возвращает  LinksStorages
func GetLinksStorage() linksstorage.LinksStorages {
	return GlobalContainer.Get("linksstorage").(linksstorage.LinksStorages)
}
