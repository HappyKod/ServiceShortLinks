package mem_links_storage

import (
	"HappyKod/ServiceShortLinks/utils"
	"sync"
)

type connect struct {
	mu    sync.Mutex
	cache map[string]string
}

type MemLinksStorage struct {
	Connect *connect
}

// New инициализация хранилища
func New() (*MemLinksStorage, error) {
	return &MemLinksStorage{
		Connect: &connect{cache: make(map[string]string)},
	}, nil
}

// Ping проверка хранилища
func (MS MemLinksStorage) Ping() (bool, error) {
	return true, nil
}

// Get получаем значение по ключу
func (MS MemLinksStorage) Get(key string) (string, error) {
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	return MS.Connect.cache[key], nil
}

// Put добавляем значение по ключу
func (MS MemLinksStorage) Put(key string, url string) error {
	MS.Connect.mu.Lock()
	MS.Connect.cache[key] = url
	MS.Connect.mu.Unlock()
	return nil
}

// CreateUniqKey Создаем уникальный ключ для записи
func (MS MemLinksStorage) CreateUniqKey() (string, error) {
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
