package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage"

	"github.com/sarulabs/di"
)

// GlobalContainer DI контейнер для общего доступа к данным
var GlobalContainer di.Container

const (
	CookeSessionName = "User-session" // Наименование ключа сессии
	CookeUserIDName  = "UserID"       // Наименование ключа UserID
	CookeUserIDLen   = 16             // Срез ключа сессии
)

// GetLinksStorage возвращает  LinksStorages
func GetLinksStorage() linksstorage.LinksStorages {
	return GlobalContainer.Get("linksstorage").(linksstorage.LinksStorages)
}
