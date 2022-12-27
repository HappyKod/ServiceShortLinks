package handlers

import (
	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func ExampleGivHandler() {
	cfg := models.Config{}
	key := "test1"
	err := container.BuildContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	link := models.Link{ShortKey: key, FullURL: "https://github.com/HappyKod/ServiceShortLinks"}
	err = constans.GetLinksStorage().PutShortLink(key, link)
	if err != nil {
		log.Fatal(err)
	}
	router := Router()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
	router.ServeHTTP(w, req)
}

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
			link := models.Link{ShortKey: tt.keyInit, FullURL: tt.want.responseLocation}
			assert.Equal(t, constans.GetLinksStorage().PutShortLink(tt.keyInit, link), nil)
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
				if err != nil && errors.Is(os.ErrNotExist, err) {
					t.Fatal(err)
				}
			}
		})

	}
}
