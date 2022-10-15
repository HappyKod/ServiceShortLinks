package users_storage

type UsersStorage interface {
	Ping() (bool, error)
	Put(key string, link string) error
	Get(key string) ([]string, error)
}
