package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/models"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func ExamplePutHandler() {
	cfg := models.Config{}
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://github.com/HappyKod/ServiceShortLinks")))
	router.ServeHTTP(w, req)
}

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
			requestBody:   "https://yandex.ru/",
			want: want{
				responseCode: http.StatusConflict,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test4",
			requestPath:   "/",
			requestMethod: http.MethodPost,
			requestBody:   "InvalidUrl",
			want: want{
				responseCode: http.StatusBadRequest,
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
		})
	}
}
