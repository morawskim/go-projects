package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func convertBufToImage(imgBuf []byte, imgType string) (image.Image, error) {
	var img image.Image
	var err error
	const webpMax = 16383

	if imgType == "jpg" {
		img, err = jpeg.Decode(bytes.NewReader(imgBuf))
	} else if imgType == "png" {
		img, err = png.Decode(bytes.NewReader(imgBuf))
	}

	if err != nil || img == nil {
		return nil, fmt.Errorf("image bufor is corrupted: %v", err)
	}

	x, y := img.Bounds().Max.X, img.Bounds().Max.Y
	if x > webpMax || y > webpMax {
		return nil, fmt.Errorf("image is too large. Max %d", webpMax)
	}

	return img, nil
}

func isSupportedMimeType(mimeType string) (bool, string) {
	switch mimeType {
	case
		"image/jpeg",
		"image/jpg":
		return true, "jpg"
	case "image/png":
		return true, "png"
	}

	return false, ""
}

func main() {
	var webpBuf, webpThumbnailBuf bytes.Buffer
	var err error

	filePathFromArgs := flag.String("image", "", "path to image")
	flag.Parse()

	absFilePath, err := filepath.Abs(*filePathFromArgs)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot parse path %v", *filePathFromArgs)
		os.Exit(1)
	}

	filePath := &absFilePath

	if !fileExists(*filePath) {
		_, _ = fmt.Fprintf(os.Stderr, "File %v not exists\n", *filePath)
		os.Exit(1)
	}

	mtype, err := mimetype.DetectFile(*filePath)
	ok, imgFormat := isSupportedMimeType(mtype.String())
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "File %v does not look like image (%v)\n", *filePath, mtype.String())
		os.Exit(1)
	}

	imgFileBuf, err := os.ReadFile(*filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "File %v cannot be read - %v\n", *filePath, err)
		os.Exit(1)
	}

	img, err := convertBufToImage(imgFileBuf, imgFormat)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot parse image %v for reason %v\n", *filePath, err)
		os.Exit(1)
	}

	// convert to webp and create miniature
	thumbnail := imaging.Resize(img, 100, 0, imaging.Lanczos)

	if err = webp.Encode(&webpBuf, img, &webp.Options{}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot encode image to webp format - %v", err)
		os.Exit(1)
	}

	if err = webp.Encode(&webpThumbnailBuf, thumbnail, &webp.Options{}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot create thumbnail - %v", err)
		os.Exit(1)
	}

	webpOutputPath := filepath.Join(
		filepath.Dir(absFilePath),
		strings.TrimSuffix(filepath.Base(absFilePath), filepath.Ext(absFilePath))+".webp",
	)

	webpThumbnailOutputPath := filepath.Join(
		filepath.Dir(absFilePath),
		strings.TrimSuffix(filepath.Base(absFilePath), filepath.Ext(absFilePath))+"_thumbnail.webp",
	)

	if err = os.WriteFile(webpOutputPath, webpBuf.Bytes(), 0644); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot save webp - %v", err)
		os.Exit(1)
	}

	if err = os.WriteFile(webpThumbnailOutputPath, webpThumbnailBuf.Bytes(), 0644); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot save thumbnail - %v", err)
		os.Exit(1)
	}
}
