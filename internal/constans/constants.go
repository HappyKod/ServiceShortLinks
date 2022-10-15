package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/links-storage"
	users_storage "HappyKod/ServiceShortLinks/internal/storage/users-storage"
	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

const (
	CookeSessionName = "User-session"
	CookeUserIDName  = "UserID"
)

func GetLinksStorage() links_storage.LinksStorages {
	return GlobalContainer.Get("links-storage").(links_storage.LinksStorages)
}
func GetUsersStorage() users_storage.UsersStorage {
	return GlobalContainer.Get("users-storage").(users_storage.UsersStorage)
}
