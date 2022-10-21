package linksstorage

import "HappyKod/ServiceShortLinks/internal/models"

// LinksStorages Набор методов для работы с хранилищем
type LinksStorages interface {
	Ping() error
	PutShortLink(key string, link models.Link) error
	ManyPutShortLink(links []models.Link) error
	GetShortLink(key string) (string, error)
	GetKey(fullURL string) (string, error)
	GetShortLinkUser(UserID string) ([]models.Link, error)
}
