package FilesLinksStorage

import (
	"HappyKod/ServiceShortLinks/utils"
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type connect struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	mu      *sync.Mutex
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
			mu:      new(sync.Mutex),
		},
		FileNAME: FileNAME,
	}, nil
}

// Ping проверка хранилища
func (FS FileLinksStorage) Ping() (bool, error) {
	_, err := FS.Connect.file.Stat()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Get получаем значение по ключу
func (FS FileLinksStorage) Get(key string) (string, error) {
	FS.Connect.mu.Lock()
	defer FS.Connect.mu.Unlock()
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
		if event[key] != "" {
			return event[key], nil
		}
	}
	return "", nil
}

// Put добавляем значение по ключу
func (FS FileLinksStorage) Put(key string, url string) error {
	FS.Connect.mu.Lock()
	defer FS.Connect.mu.Unlock()
	structMAP := map[string]string{key: url}
	err := FS.Connect.encoder.Encode(&structMAP)
	if err != nil {
		return err
	}
	return nil
}

// CreateUniqKey Создаем уникальный ключ для записи
func (FS FileLinksStorage) CreateUniqKey() (string, error) {
	var key string
	for {
		key = utils.GeneratorStringUUID()
		url, err := FS.Get(key)
		if err != nil {
			return "", err
		}
		if url == "" {
			break
		}
	}
	return key, nil
}

// Close закрываем соединение (файл)
func (FS FileLinksStorage) Close() error {
	err := FS.Connect.file.Close()
	if err != nil {
		return err
	}
	return nil
}
