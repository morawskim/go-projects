package main

import (
	"bytes"
	"errors"
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

	if errors.Is(err, os.ErrNotExist) {
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

func getAbsolutePathToFile(filePathFromArgs *string) (string, error) {
	absFilePath, err := filepath.Abs(*filePathFromArgs)
	if err != nil {
		return "", fmt.Errorf("cannot parse path %v", *filePathFromArgs)
	}

	if !fileExists(absFilePath) {
		return "", fmt.Errorf("file %v not exists\n", absFilePath)
	}

	return absFilePath, nil
}

func convertImgToWebp(img image.Image) (*bytes.Buffer, error) {
	var webpBuf bytes.Buffer

	if err := webp.Encode(&webpBuf, img, &webp.Options{}); err != nil {
		return nil, fmt.Errorf("cannot encode image to webp format - %v", err)
	}

	return &webpBuf, nil
}

func createWebpThumbnail(img image.Image) (*bytes.Buffer, error) {
	var webpThumbnailBuf bytes.Buffer
	thumbnail := imaging.Resize(img, 100, 0, imaging.Lanczos)

	if err := webp.Encode(&webpThumbnailBuf, thumbnail, &webp.Options{}); err != nil {
		return nil, fmt.Errorf("cannot create thumbnail - %v", err)
	}

	return &webpThumbnailBuf, nil
}

func saveWebpImg(path string, buffer *bytes.Buffer) error {
	if err := os.WriteFile(path, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("cannot save webp - %v", err)
	}

	return nil
}

func main() {
	var err error

	filePathFromArgs := flag.String("image", "", "path to image")
	flag.Parse()

	filePath, err := getAbsolutePathToFile(filePathFromArgs)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot get abosulte path for file %v: %v", *filePathFromArgs, err)
		os.Exit(1)
	}

	mtype, err := mimetype.DetectFile(filePath)
	ok, imgFormat := isSupportedMimeType(mtype.String())
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "File %v does not look like image (%v)\n", filePath, mtype.String())
		os.Exit(1)
	}

	imgFileBuf, err := os.ReadFile(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "File %v cannot be read - %v\n", filePath, err)
		os.Exit(1)
	}

	img, err := convertBufToImage(imgFileBuf, imgFormat)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot parse image %v for reason %v\n", filePath, err)
		os.Exit(1)
	}

	// convert to webp and create miniature
	webpBuf, err := convertImgToWebp(img)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot encode image to webp format - %v", err)
		os.Exit(1)
	}

	webpThumbnailBuf, err := createWebpThumbnail(img)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot create thumbnail - %v", err)
		os.Exit(1)
	}

	err = saveWebpImg(
		filepath.Join(
			filepath.Dir(filePath),
			strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".webp",
		),
		webpBuf,
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot save webp - %v", err)
		os.Exit(1)
	}

	err = saveWebpImg(
		filepath.Join(
			filepath.Dir(filePath),
			strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+"_thumbnail.webp",
		),
		webpThumbnailBuf,
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot save thumbnail - %v", err)
		os.Exit(1)
	}
}
