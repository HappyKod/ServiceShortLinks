package links_storage

// LinksStorages Набор методов для работы с хранилищем
type LinksStorages interface {
	Ping() (bool, error)
	Put(key string, url string) error
	Get(key string) (string, error)
	CreateUniqKey() (string, error)
}
