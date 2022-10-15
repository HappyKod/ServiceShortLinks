package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/models"
	"bytes"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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
