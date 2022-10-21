package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage"
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
