package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/models"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

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
		{
			name:          "Получаем ошибку при генерации ссылки test4 Mem хранилище",
			requestPath:   "/api/shorten",
			requestMethod: http.MethodPost,
			requestBody:   `{"url": "https://yandex.ru"}`,
			want: want{
				responseCode: http.StatusConflict,
			},
		},
	}
	err := container.BuildContainer(models.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
