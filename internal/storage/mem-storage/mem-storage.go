package mem_storage

import "sync"

type connect struct {
	mu    sync.Mutex
	cache map[string]string
}

type MemStorage struct {
	Connect *connect
}

//Init иницилизация хранилища
func (MS MemStorage) Init() error {
	MS.Connect = &connect{
		cache: make(map[string]string),
	}
	return nil
}

//Ping проверка харнилища
func (MS MemStorage) Ping() (bool, error) {
	return true, nil
}

//Get получаем значение по ключу
func (MS MemStorage) Get(key string) (string, error) {
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	return MS.Connect.cache[key], nil
}

//Put добавляем занчение по ключу
func (MS MemStorage) Put(key string, url string) error {
	MS.Connect.mu.Lock()
	MS.Connect.cache[key] = url
	MS.Connect.mu.Unlock()
	return nil
}
