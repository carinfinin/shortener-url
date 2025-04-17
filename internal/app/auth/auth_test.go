package auth

import "testing"

func BenchmarkEncodeToken(b *testing.B) {
	token := GenerateToken()
	for i := 0; i < b.N; i++ {
		EncodeToken(token)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken()
	}
}
