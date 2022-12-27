package utils

import "testing"

func BenchmarkValidatorURL(b *testing.B) {
	text := "https://github.com/HappyKod/ServiceShortLinks"
	b.Run(b.Name(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValidatorURL(text)
		}
	})
}

func BenchmarkGeneratorStringUUID(b *testing.B) {
	b.ReportAllocs()
	b.Run(b.Name(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GeneratorStringUUID()
		}
	})
}
