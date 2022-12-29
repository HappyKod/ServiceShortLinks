package utils

import "testing"

func BenchmarkValidatorURL(b *testing.B) {
	text := "https://github.com/HappyKod/ServiceShortLinks"
	b.Run(b.Name(), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ValidatorURL(text)
		}
	})
}

func BenchmarkGeneratorStringUUID(b *testing.B) {
	b.Run(b.Name(), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			GeneratorStringUUID()
		}
	})
}

func TestValidatorURL(t *testing.T) {
	type args struct {
		rawText string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{"https://yandex.ru/"},
			want: true,
		},
		{
			name: "invalid",
			args: args{"yandex.ru/"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatorURL(tt.args.rawText); got != tt.want {
				t.Errorf("ValidatorURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
