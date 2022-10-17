package linksstorage

// LinksStorages Набор методов для работы с хранилищем
type LinksStorages interface {
	Ping() error
	PutShortLink(key string, url string) error
	ManyPutShortLink(urls []string) (map[string]string, error)
	GetShortLink(key string) (string, error)
	GetKey(fullURL string) (string, error)
	CreateUniqKey() (string, error)
}
