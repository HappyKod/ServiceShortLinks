package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage"

	"github.com/sarulabs/di"
)

// GlobalContainer DI контейнер
var GlobalContainer di.Container

const (
	//CookeSessionName наименование ключа сессии
	CookeSessionName = "User-session"
	//CookeUserIDName наименование ключа UserID
	CookeUserIDName = "UserID"
	//CookeUserIDLen срез ключа сессии
	CookeUserIDLen = 16
)

// GetLinksStorage возвращает  LinksStorages
func GetLinksStorage() linksstorage.LinksStorages {
	return GlobalContainer.Get("linksstorage").(linksstorage.LinksStorages)
}
