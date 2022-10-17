package linksstorage

// LinksStorages Набор методов для работы с хранилищем
type LinksStorages interface {
	Ping() error
	Put(key string, url string) error
	ManyPut(urls []string) (map[string]string, error)
	Get(key string) (string, error)
	CreateUniqKey() (string, error)
}
