package container

import (
	"testing"

	"github.com/HappyKod/ServiceShortLinks/internal/models"
)

func TestBuildContainer(t *testing.T) {
	type args struct {
		cfg models.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Запуск контейнера с локальным хранилищем",
			args: args{
				cfg: models.Config{
					Address:         "localhost:8080",
					BaseURL:         "http://localhost:8080",
					FileStoragePATH: "",
				},
			},
		},
		{
			name: "Запуск контейнера с файловым хранилищем",
			args: args{
				cfg: models.Config{
					Address:         "localhost:8080",
					BaseURL:         "http://localhost:8080",
					FileStoragePATH: "file.json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BuildContainer(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("BuildContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
