package memusersstorage

import (
	"reflect"
	"testing"
)

func TestMemUsersStorage(t *testing.T) {
	type args struct {
		key  string
		link string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{
			name: "загрузка user1 ссылки",
			args: args{
				link: "test1",
				key:  "user1",
			},
			want: []string{"test1"},
		},
		{
			name: "загрузка user2 ссылки",
			args: args{
				link: "test2",
				key:  "user2",
			},
			want: []string{"test2"},
		},
		{
			name: "загрузка user1 ссылки",
			args: args{
				link: "test3",
				key:  "user1",
			},
			want: []string{"test1", "test3"},
		},
	}
	MS, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MS.Put(tt.args.key, tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("PutShortLink() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := MS.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutShortLink() error = %v, wantErr %v", err, tt.wantErr)

			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShortLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
