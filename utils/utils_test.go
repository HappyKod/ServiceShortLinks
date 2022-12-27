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
