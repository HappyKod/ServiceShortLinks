package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestGivHandler тест метода GivHandler
func TestGivHandler(t *testing.T) {
	type want struct {
		responseCode     int
		responseLocation string
	}
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		keyInit       string
		key           string
		want          want
		cfg           models.Config
	}{
		{
			name:          "Получаем ссылку по ключу test1 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test1",
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://github.com/HappyKod/ServiceShortLinks",
			},
		},
		{
			name:          "Получаем ссылку по ключу test2 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test2",
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://yandex.ru/",
			},
		},
		{
			name:          "Получаем ссылку по не верному key test2 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test3",
			key:           "test3Invalid",
			want: want{
				responseCode:     http.StatusBadRequest,
				responseLocation: "https://yandex.ru/",
			},
		},
		{
			name:          "Получаем ссылку по ключу test1 File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test1",
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://github.com/HappyKod/ServiceShortLinks",
			},
		},
		{
			name:          "Получаем ссылку по ключу test2 File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test2",
			cfg:           models.Config{FileStoragePATH: "test2.json"},
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://yandex.ru/",
			},
		},
		{
			name:          "Получаем ссылку по не верному key test2 File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test3",
			key:           "test3Invalid",
			cfg:           models.Config{FileStoragePATH: "test3.json"},
			want: want{
				responseCode:     http.StatusBadRequest,
				responseLocation: "https://yandex.ru/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := container.BuildContainer(tt.cfg)
			if err != nil {
				t.Fatal(err)
			}
			if tt.key == "" {
				tt.key = tt.keyInit
			}
			//Наполняем тестовыми данными
			assert.Equal(t, constans.GlobalContainer.Get("links-storage").(storage.Storages).Put(tt.keyInit, tt.want.responseLocation), nil)
			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath+tt.key, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			if w.Code == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.responseLocation, w.Header().Get("Location"))
			} else {
				assert.NotEqual(t, tt.want.responseLocation, w.Header().Get("Location"))
			}
			if tt.cfg.FileStoragePATH != "" {
				err = os.Remove(tt.cfg.FileStoragePATH)
				if err != nil {
					t.Fatal(err)
				}
			}
		})

	}
}

// TestPutHandler тест метода PutHandler
func TestPutHandler(t *testing.T) {
	type want struct {
		responseCode int
	}
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		want          want
		cfg           models.Config
	}{
		{
			name:          "Генерируем сокращенную ссылку test1 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://github.com/HappyKod/ServiceShortLinks",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru/",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3 Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "InvalidUrl",
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test1 FIle хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://github.com/HappyKod/ServiceShortLinks",
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2 File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru/",
			cfg:           models.Config{FileStoragePATH: "test2.json"},
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3 File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "InvalidUrl",
			cfg:           models.Config{FileStoragePATH: "test3.json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := container.BuildContainer(tt.cfg)
			if err != nil {
				t.Fatal(err)
			}
			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			if tt.cfg.FileStoragePATH != "" {
				err := os.Remove(tt.cfg.FileStoragePATH)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestGIVPUT Тестируем связку PutHandler GivHandler
func TestGIVGET(t *testing.T) {
	type want struct {
		responseCode int
	}
	var Key string
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		want          want
		cfg           models.Config
	}{
		{
			name:          "Кодируем ссылку Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу Mem хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name:          "Кодируем ссылку FIle хранилище",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru",
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу File хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
	}
	err := container.BuildContainer(models.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.FileStoragePATH != "" {
				err := container.BuildContainer(tt.cfg)
				if err != nil {
					t.Fatal(err)
				}
			}
			router := Router()
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.requestMethod == http.MethodGet {
				fmt.Println(Key)
				req = httptest.NewRequest(http.MethodGet, tt.requestPath+Key, nil)
			} else if tt.requestMethod == http.MethodPost {
				req = httptest.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			}
			router.ServeHTTP(w, req)
			fmt.Println(w.Body.String())
			assert.Equal(t, tt.want.responseCode, w.Code)
			Key = w.Body.String()
		})
	}
	for _, v := range tests {
		if v.cfg.FileStoragePATH != "" {
			err = os.Remove(v.cfg.FileStoragePATH)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// TestPutHandler тест метода PutAPIHandler
func TestPutApiHandler(t *testing.T) {
	type want struct {
		responseCode int
	}
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		want          want
		cfg           models.Config
	}{
		{
			name:          "Генерируем сокращенную ссылку test1 Mem хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://github.com/HappyKod/ServiceShortLinks"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2 Mem хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://yandex.ru"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3 Mem хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "InvalidUrl"}`,
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := container.BuildContainer(tt.cfg)
			if err != nil {
				t.Fatal(err)
			}
			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			if tt.cfg.FileStoragePATH != "" {
				err = os.Remove(tt.cfg.FileStoragePATH)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestPutApiGET Тестируем связку GivHandler PutAPIHandler
func TestPutApiGET(t *testing.T) {
	type want struct {
		responseCode int
	}
	var Key string
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		want          want
		cfg           models.Config
	}{
		{
			name:          "Кодируем ссылку Mem Хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://github.com/HappyKod/ServiceShortLinks"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу Mem Хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name:          "Кодируем ссылку File Хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://github.com/HappyKod/ServiceShortLinks"}`,
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу File Хранилище",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			cfg:           models.Config{FileStoragePATH: "test1.json"},
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
	}
	err := container.BuildContainer(models.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.FileStoragePATH != "" {
				err := container.BuildContainer(tt.cfg)
				if err != nil {
					t.Fatal(err)
				}
			}
			router := Router()
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.requestMethod == http.MethodGet {
				req = httptest.NewRequest(http.MethodGet, tt.requestPath+Key, nil)
				router.ServeHTTP(w, req)
				assert.Equal(t, tt.want.responseCode, w.Code)
			} else if tt.requestMethod == http.MethodPost {
				req = httptest.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
				router.ServeHTTP(w, req)
				assert.Equal(t, tt.want.responseCode, w.Code)
				keyMAP := make(map[string]string)
				err = json.Unmarshal(w.Body.Bytes(), &keyMAP)
				if err != nil {
					t.Fatal(err)
				}
				Key = keyMAP["result"]
			}
		})
	}
	for _, v := range tests {
		if v.cfg.FileStoragePATH != "" {
			err = os.Remove(v.cfg.FileStoragePATH)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
