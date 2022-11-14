package util

import (
	"errors"
	"testing"
)

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, e := HashPassword("123456")
		println(b)
		if e != nil {
			panic(e)
		}
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !CheckPasswordHash("123456", "$2a$14$P.sU6GpIr7gcyXzcsR/ZZOBs5powMyX1/x.y1uoqwBu6vPFemtNwO") {
			panic(errors.New("check fail"))
		}
	}
}
