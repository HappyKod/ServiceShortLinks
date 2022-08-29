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

func TestH(t *testing.T) {
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
