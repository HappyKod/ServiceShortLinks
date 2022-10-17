package memlinksstorage

import (
	"HappyKod/ServiceShortLinks/utils"
	"errors"
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
func (MS MemLinksStorage) Ping() error {
	return nil
}

// GetShortLink получаем значение по ключу
func (MS MemLinksStorage) GetShortLink(key string) (string, error) {
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	return MS.Connect.cache[key], nil
}

// PutShortLink добавляем значение по ключу
func (MS MemLinksStorage) PutShortLink(key string, url string) error {
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

// ManyPutShortLink добавляем множества значений
func (MS MemLinksStorage) ManyPutShortLink(urls []string) (map[string]string, error) {
	shortURLS := make(map[string]string)
	for _, url := range urls {
		key, err := MS.CreateUniqKey()
		if err != nil {
			return nil, err
		}
		if err = MS.PutShortLink(key, url); err != nil {
			return nil, err
		}
		MS.Connect.mu.Lock()
		shortURLS[key] = url
		MS.Connect.mu.Unlock()
	}
	return shortURLS, nil
}

func (MS MemLinksStorage) GetKey(fullURL string) (string, error) {
	MS.Connect.mu.Lock()
	localCache := MS.Connect.cache
	MS.Connect.mu.Unlock()
	for k, v := range localCache {
		if v == fullURL {
			return k, nil
		}
	}
	return "", errors.New("ссылка не найдена")
}
