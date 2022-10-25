package memlinksstorage

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"errors"
	"sync"
)

type MemLinksStorage struct {
	mu    *sync.RWMutex
	cache map[string]models.Link
}

// New инициализация хранилища
func New() (*MemLinksStorage, error) {
	return &MemLinksStorage{
		cache: make(map[string]models.Link),
		mu:    new(sync.RWMutex)}, nil
}

// Ping проверка хранилища
func (MS MemLinksStorage) Ping() error {
	return nil
}

// GetShortLink получаем полную ссылку по ключу
func (MS MemLinksStorage) GetShortLink(key string) (string, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	return MS.cache[key].FullURL, nil
}

// PutShortLink добавляем models.Link по ключу
func (MS MemLinksStorage) PutShortLink(key string, link models.Link) error {
	_, err := MS.GetKey(link.FullURL)
	if !errors.Is(err, constans.ErrorNotFindFullURL) {
		return constans.ErrorNoUNIQUEFullURL
	}
	MS.mu.Lock()
	MS.cache[key] = link
	MS.mu.Unlock()
	return nil
}

// ManyPutShortLink добавляем множества models.Link
func (MS MemLinksStorage) ManyPutShortLink(links []models.Link) error {
	for _, link := range links {
		if err := MS.PutShortLink(link.ShortKey, link); err != nil {
			return err
		}
	}
	return nil
}

// GetKey получаем значение ключа по полной ссылке
func (MS MemLinksStorage) GetKey(fullURL string) (string, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for k, link := range MS.cache {
		if link.FullURL == fullURL {
			return k, nil
		}
	}
	return "", constans.ErrorNotFindFullURL
}

// GetShortLinkUser получаем все models.Link который добавил пользователь
func (MS MemLinksStorage) GetShortLinkUser(UserID string) ([]models.Link, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var linksUser []models.Link
	for _, link := range MS.cache {
		if link.UserID == UserID {
			linksUser = append(linksUser, link)
		}
	}
	return linksUser, nil
}
