package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/HappyKod/ServiceShortLinks/internal/app/container"
	"github.com/HappyKod/ServiceShortLinks/internal/models"
)

func TestPutAPIBatchHandler(t *testing.T) {
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
			requestPath:   "/api/shorten/batch",
			requestMethod: http.MethodPost,
			requestBody: `[{
								"correlation_id": "11111",
								"original_url": "https://github.com/HappyKod/ServiceShortLinks"
							}]`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Генерируем сокращенную ссылку test2 Mem хранилище",
			requestPath:   "/api/shorten/batch",
			requestMethod: http.MethodPost,
			requestBody: `[
							{
								"correlation_id": "11111",
								"original_url": "https://www.twitch.tv/videos/16265634682"
							},
							{
								"correlation_id": "2222",
								"original_url": "https://www.twitch.tv/videos/162656346812"
							}
						]`,
			want: want{
				responseCode: http.StatusCreated,
			},
		},
		{
			name:          "Получаем ошибку при генерации ссылки test3 Mem хранилище",
			requestPath:   "/api/shorten/batch",
			requestMethod: http.MethodPost,
			requestBody: `[{
								"correlation_id": "11111",
								"original_url": "Invalid"
							}]`,
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
			if tt.cfg.FileStoragePATH != "" {
				err = os.Remove(tt.cfg.FileStoragePATH)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
