package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func BenchmarkConvertToWebp(b *testing.B) {
	imgFileBuf, err := os.ReadFile("3000.jpg")
	assert.NoError(b, err)

	img, err := convertBufToImage(imgFileBuf, "jpg")
	assert.NoError(b, err)

	b.ResetTimer()

	b.Run("quality: 10%", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := convertImgToWebp(img, 10.0)
			assert.NoError(b, err)
		}
	})

	b.Run("quality: 50%", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := convertImgToWebp(img, 50.0)
			assert.NoError(b, err)
		}
	})

	b.Run("quality: 90%", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := convertImgToWebp(img, 90.0)
			assert.NoError(b, err)
		}
	})
}
