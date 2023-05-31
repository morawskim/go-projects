package main

import "testing"

func BenchmarkPassGoStringToC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		passGoStringToC()
	}
}

func BenchmarkPassStructToC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		passStructToC()
	}
}

func BenchmarkPossibleMemoryLeak(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr := possibleMemoryLeak()
		freeMemoryChunk(ptr)
	}
}
