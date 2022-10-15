package linksstorage

// LinksStorages Набор методов для работы с хранилищем
type LinksStorages interface {
	Ping() error
	Put(key string, url string) error
	Get(key string) (string, error)
	CreateUniqKey() (string, error)
}
