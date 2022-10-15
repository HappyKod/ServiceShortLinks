package usersstorage

type UsersStorage interface {
	Ping() error
	Put(key string, link string) error
	Get(key string) ([]string, error)
}
