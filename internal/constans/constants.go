package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/usersstorage"
	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

const (
	CookeSessionName = "User-session"
	CookeUserIDName  = "UserID"
)

func GetLinksStorage() linksstorage.LinksStorages {
	return GlobalContainer.Get("linksstorage").(linksstorage.LinksStorages)
}
func GetUsersStorage() usersstorage.UsersStorage {
	return GlobalContainer.Get("usersstorage").(usersstorage.UsersStorage)
}
