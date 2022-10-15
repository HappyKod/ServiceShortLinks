package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/Linksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/UsersStorage"
	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

const (
	CookeSessionName = "User-session"
	CookeUserIDName  = "UserID"
)

func GetLinksStorage() Linksstorage.LinksStorages {
	return GlobalContainer.Get("Linksstorage").(Linksstorage.LinksStorages)
}
func GetUsersStorage() UsersStorage.UsersStorage {
	return GlobalContainer.Get("UsersStorage").(UsersStorage.UsersStorage)
}
