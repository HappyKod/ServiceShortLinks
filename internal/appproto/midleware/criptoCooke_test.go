package midleware

import (
	"log"
	"testing"

	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/models"
)

func Test_validCookie(t *testing.T) {
	type args struct {
		cooke string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			name: "valid cookie",
			args: args{
				cooke: "39636466363139662d326137632d3439463dede60c4a8829edda83bfd46dc019add7966c028655357a1e086706157be8",
			},
			want:  "9cdf619f-2a7c-49",
			want1: true,
		},
		{
			name: "invalid cookie",
			args: args{
				cooke: "1139636466363139662d326137632d3439463dede60c4a8829edda83bfd46dc019add7966c028655357a1e086706157be8",
			},
			want1: false,
		},
	}
	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := validCookie(tt.args.cooke)
			if (err != nil) != tt.wantErr {
				t.Errorf("validCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validCookie() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("validCookie() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
