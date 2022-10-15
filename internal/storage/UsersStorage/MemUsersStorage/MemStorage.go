package MemUsersStorage

import (
	"sync"
)

type connect struct {
	mu    sync.Mutex
	cache map[string][]string
}

type MemUsersStorage struct {
	Connect *connect
}

// New инициализация хранилища
func New() (*MemUsersStorage, error) {
	return &MemUsersStorage{
		Connect: &connect{cache: make(map[string][]string)},
	}, nil
}

// Ping проверка хранилища
func (MS MemUsersStorage) Ping() (bool, error) {
	return true, nil
}

// Get получаем значение по ключу
func (MS MemUsersStorage) Get(key string) ([]string, error) {
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	return MS.Connect.cache[key], nil
}

// Put добавляем значение по ключу
func (MS MemUsersStorage) Put(key string, url string) error {
	links, err := MS.Get(key)
	if err != nil {
		return err
	}
	MS.Connect.mu.Lock()
	defer MS.Connect.mu.Unlock()
	if len(links) == 0 {
		MS.Connect.cache[key] = []string{url}
	} else {
		links = append(links, url)
		MS.Connect.cache[key] = links
	}
	return nil
}
