package storage

//Storages
//Набор методов для работы с хранилищем
type Storages interface {
	Ping() (bool, error)
	Put(key string, url string) error
	Get(key string) (string, error)
}
