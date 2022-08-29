package handlers

import (
	"ServiceShortLinks/internal/constans"
	mem_storage "ServiceShortLinks/internal/storage/memstorage"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Key string

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
				responseCode:     307,
				responseLocation: "https://github.com/HappyKod/ServiceShortLinks",
			},
		},
		{
			name:          "Получаем ссылку по ключу test2",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			keyInit:       "test2",
			want: want{
				responseCode:     307,
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
				responseCode:     400,
				responseLocation: "https://yandex.ru/",
			},
		},
	}
	storage, err := mem_storage.Init()
	if err != nil {
		log.Fatalln(errors.New("Ошибка иницилизации mem_storage " + err.Error()))
	}
	//иницилизирум глобальное хранилище
	constans.GlobalStorage = mem_storage.MemStorage{Connect: storage}

	//Наполяем тестовыми данными

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key == "" {
				tt.key = tt.keyInit
			}
			//Наполняем тестовыми данными
			assert.Equal(t, constans.GlobalStorage.Put(tt.keyInit, tt.want.responseLocation), nil)
			router := Router()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath+tt.key, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			if w.Code == 307 {
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
				responseCode: 201,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "https://yandex.ru/",
			want: want{
				responseCode: 201,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "InvalidUrl",
			want: want{
				responseCode: 400,
			},
		},
	}
	storage, err := mem_storage.Init()
	if err != nil {
		log.Fatalln(errors.New("Ошибка иницилизации mem_storage " + err.Error()))
	}
	//иницилизирум глобальное хранилище
	constans.GlobalStorage = mem_storage.MemStorage{Connect: storage}
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
				responseCode: 201,
			},
		},
		{
			name:          "Получаем ссылку по ключу",
			requestPath:   "/",
			requestMethod: http.MethodGet,
			want: want{
				responseCode: 307,
			},
		},
	}
	storage, err := mem_storage.Init()
	if err != nil {
		log.Fatalln(errors.New("Ошибка иницилизации mem_storage " + err.Error()))
	}
	//иницилизирум глобальное хранилище
	constans.GlobalStorage = mem_storage.MemStorage{Connect: storage}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Router()
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.requestMethod == http.MethodGet {
				fmt.Println(tt.requestPath + Key)
				req = httptest.NewRequest(http.MethodGet, Key, nil)
			} else if tt.requestMethod == http.MethodPost {
				req = httptest.NewRequest(http.MethodPost, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.responseCode, w.Code)
			Key = w.Body.String()
		})
	}
	fmt.Println(Key)
}
