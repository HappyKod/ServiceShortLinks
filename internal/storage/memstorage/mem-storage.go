package memstorage

import (
	"HappyKod/ServiceShortLinks/utils"
	"sync"
)

type connect struct {
	mu    sync.Mutex
	cache map[string]string
}

type MemStorage struct {
	Connect *connect
}

// New инициализация хранилища
func New() (*MemStorage, error) {
	return &MemStorage{
		Connect: &connect{cache: make(map[string]string)},
	}, nil
}

// Ping проверка хранилища
func (MS MemStorage) Ping() (bool, error) {
	return true, nil
}

// Get получаем значение по ключу
func (MS MemStorage) Get(key string) (string, error) {
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	return MS.Connect.cache[key], nil
}

// Put добавляем значение по ключу
func (MS MemStorage) Put(key string, url string) error {
	MS.Connect.mu.Lock()
	MS.Connect.cache[key] = url
	MS.Connect.mu.Unlock()
	return nil
}

// CreateUniqKey Создаем уникальный ключ для записи
func (MS MemStorage) CreateUniqKey() (string, error) {
	var key string
	var url string
	for {
		key = utils.GeneratorStringUUID()
		MS.Connect.mu.Lock()
		MS.Connect.cache[key] = url
		MS.Connect.mu.Unlock()
		if url == "" {
			break
		}
	}
	return key, nil
}
