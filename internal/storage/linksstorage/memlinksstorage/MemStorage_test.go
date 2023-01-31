package memlinksstorage

import (
	"sync"
	"testing"

	"github.com/HappyKod/ServiceShortLinks/internal/models"
)

func TestMemLinksStorage_GetShortLink(t *testing.T) {
	type fields struct {
		mu    *sync.RWMutex
		cache map[string]models.Link
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Link
		wantErr bool
	}{
		{
			name: "GetShortLink",
			fields: fields{
				mu:    &sync.RWMutex{},
				cache: make(map[string]models.Link),
			},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MS, err := New()
			if err != nil {
				t.Fatal(err)
			}
			_, err = MS.GetShortLink(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMemLinksStorage_ManyPutShortLink(t *testing.T) {
	type fields struct {
		mu    *sync.RWMutex
		cache map[string]models.Link
	}
	type args struct {
		links []models.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "GetShortLink",
			fields: fields{
				mu:    &sync.RWMutex{},
				cache: make(map[string]models.Link),
			},
			args: args{
				links: []models.Link{
					models.Link{
						FullURL: "https://www.google.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MS, err := New()
			if err != nil {
				t.Fatal(err)
			}
			if err := MS.ManyPutShortLink(tt.args.links); (err != nil) != tt.wantErr {
				t.Errorf("ManyPutShortLink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemLinksStorage_PutShortLink(t *testing.T) {
	type fields struct {
		mu    *sync.RWMutex
		cache map[string]models.Link
	}
	type args struct {
		key  string
		link models.Link
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "TestMemLinksStorage_PutShortLink", args: args{key: "key"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MS, err := New()
			if err != nil {
				t.Fatal(err)
			}
			if err = MS.PutShortLink(tt.args.key, tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("PutShortLink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *MemLinksStorage
		wantErr bool
	}{
		{
			name: "TestNew",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
