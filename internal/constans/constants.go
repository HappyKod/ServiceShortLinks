package constans

import (
	"HappyKod/ServiceShortLinks/internal/storage"
	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

func GetLinkStorage() storage.Storages {
	return GlobalContainer.Get("links-storage").(storage.Storages)
}
