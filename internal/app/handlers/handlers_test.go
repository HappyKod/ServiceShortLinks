package handlers

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage"
	"HappyKod/ServiceShortLinks/internal/storage/memstorage"
	"bytes"
	"encoding/json"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/assert/v2"
	"github.com/sarulabs/di"
	"net/http"
	"net/http/httptest"
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
	}{
		{
			name:          "Получаем ссылку по ключу test1",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test1",
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://github.com/HappyKod/ServiceShortLinks",
			},
		},
		{
			name:          "Получаем ссылку по ключу test2",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test2",
			want: want{
				responseCode:     http.StatusTemporaryRedirect,
				responseLocation: "https://yandex.ru/",
			},
		},
		{
			name:          "Получаем ссылку по не верному key test2",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test3",
			key:           "test3Invalid",
			want: want{
				responseCode:     http.StatusBadRequest,
				responseLocation: "https://yandex.ru/",
			},
		},
	}
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		t.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			linksStorage, err := memstorage.New()
			if err != nil {
				t.Fatal("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: linksStorage}, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	}{
		{
			name:          "Генерируем сокращенную ссылку test1",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://github.com/HappyKod/ServiceShortLinks",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru/",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "InvalidUrl",
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		t.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			linksStorage, err := memstorage.New()
			if err != nil {
				t.Fatal("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: linksStorage}, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
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
	}{
		{
			name:          "Кодируем ссылку",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru",
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
	}
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		t.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			linksStorage, err := memstorage.New()
			if err != nil {
				t.Fatal("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: linksStorage}, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Router()
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.requestMethod == http.MethodGet {
				req = httptest.NewRequest(http.MethodGet, Key, nil)
			} else if tt.requestMethod == http.MethodPost {
				req = httptest.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			Key = w.Body.String()
		})
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
	}{
		{
			name:          "Генерируем сокращенную ссылку test1",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://github.com/HappyKod/ServiceShortLinks"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://yandex.ru"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "InvalidUrl"}`,
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
	}
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		t.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			linksStorage, err := memstorage.New()
			if err != nil {
				t.Fatal("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: linksStorage}, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
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
	}{
		{
			name:          "Кодируем ссылку",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://github.com/HappyKod/ServiceShortLinks"}`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ссылку по ключу",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				responseCode: http.StatusTemporaryRedirect,
			},
		},
	}
	var cfg models.Config
	err := env.Parse(&cfg)
	if err != nil {
		t.Fatal(err, "Ошибка считывания конфига")
	}
	builder, _ := di.NewBuilder()
	err = builder.Add(di.Def{
		Name: "links-storage",
		Build: func(ctn di.Container) (interface{}, error) {
			linksStorage, err := memstorage.New()
			if err != nil {
				t.Fatal("Ошибка иницилизации mem_storage ", err)
			}
			return memstorage.MemStorage{Connect: linksStorage}, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	err = builder.Add(di.Def{
		Name: "server-config",
		Build: func(ctn di.Container) (interface{}, error) {
			return cfg, nil
		}})
	if err != nil {
		t.Fatal("Ошибка иницилизации контейнера", err)
	}
	constans.GlobalContainer = builder.Build()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Router()
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.requestMethod == http.MethodGet {
				req = httptest.NewRequest(http.MethodGet, Key, nil)
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
				Key = keyMAP["Result"]
			}
		})
	}
}
