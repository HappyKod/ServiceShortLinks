package fileslinksstorage

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type connect struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	mu      *sync.RWMutex
}

type FileLinksStorage struct {
	Connect  *connect
	FileNAME string
}

// New инициализации хранилища
func New(FileNAME string) (*FileLinksStorage, error) {
	file, err := os.OpenFile(FileNAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &FileLinksStorage{
		Connect: &connect{
			file:    file,
			encoder: json.NewEncoder(file),
			decoder: json.NewDecoder(file),
			mu:      new(sync.RWMutex),
		},
		FileNAME: FileNAME,
	}, nil
}

// Ping проверка хранилища
func (FS FileLinksStorage) Ping() error {
	_, err := FS.Connect.file.Stat()
	if err != nil {
		return err
	}
	return nil
}

// GetShortLink получаем значение по ключу
func (FS FileLinksStorage) GetShortLink(key string) (models.Link, error) {
	FS.Connect.mu.RLock()
	defer FS.Connect.mu.RUnlock()
	file, err := os.Open(FS.FileNAME)
	if err != nil {
		return models.Link{}, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event := make(map[string]string)
		if err = json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return models.Link{}, err
		}
		if event[key] != "" {
			var link models.Link
			err = json.Unmarshal([]byte(event[key]), &link)
			if err != nil {
				return models.Link{}, err
			}
			if link.FullURL != "" {
				return link, nil
			}
		}
	}
	return models.Link{}, nil
}

// PutShortLink добавляем значение по ключу
func (FS FileLinksStorage) PutShortLink(key string, link models.Link) error {
	_, err := FS.GetKey(link.FullURL)
	if !errors.Is(err, constans.ErrorNotFindFullURL) {
		return constans.ErrorNoUNIQUEFullURL
	}
	FS.Connect.mu.Lock()
	defer FS.Connect.mu.Unlock()
	linkStr, err := json.Marshal(link)
	if err != nil {
		return err
	}
	structMAP := map[string]string{key: string(linkStr)}
	err = FS.Connect.encoder.Encode(&structMAP)
	if err != nil {
		return err
	}
	return nil
}

// Close закрываем соединение (файл)
func (FS FileLinksStorage) Close() error {
	err := FS.Connect.file.Close()
	if err != nil {
		return err
	}
	return nil
}

// ManyPutShortLink добавляем множества значений
func (FS FileLinksStorage) ManyPutShortLink(links []models.Link) error {
	for _, link := range links {
		if err := FS.PutShortLink(link.ShortKey, link); err != nil {
			return err
		}
	}
	return nil
}

func (FS FileLinksStorage) GetKey(fullURL string) (string, error) {
	FS.Connect.mu.RLock()
	defer FS.Connect.mu.RUnlock()
	file, err := os.Open(FS.FileNAME)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event := make(map[string]string)
		if err = json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return "", err
		}
		for k, v := range event {
			var link models.Link
			err = json.Unmarshal([]byte(v), &link)
			if err != nil {
				return "", err
			}
			if link.FullURL == fullURL {
				return k, nil
			}
		}
	}
	return "", constans.ErrorNotFindFullURL
}

func (FS FileLinksStorage) GetShortLinkUser(UserID string) ([]models.Link, error) {
	FS.Connect.mu.RLock()
	defer FS.Connect.mu.RUnlock()
	file, err := os.Open(FS.FileNAME)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	var linksUser []models.Link
	for scanner.Scan() {
		event := make(map[string]string)
		if err = json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, err
		}
		for _, v := range event {
			var link models.Link
			err = json.Unmarshal([]byte(v), &link)
			if err != nil {
				return nil, err
			}
			if link.UserID == UserID {
				linksUser = append(linksUser, link)
			}
		}
	}
	return linksUser, nil
}

func (FS FileLinksStorage) DeleteShortLinkUser(UserID string, keys []string) error {
	return nil
}
