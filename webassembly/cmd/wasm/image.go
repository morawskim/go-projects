package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"syscall/js"
	"time"
)

type imageResizer struct {
	inBuf       []uint8
	sourceImg   image.Image
	onImgLoadCb js.Func
	console     js.Value
}

func newImageResizer() *imageResizer {
	return &imageResizer{
		inBuf:     nil,
		sourceImg: nil,
		console:   js.Global().Get("console"),
	}
}

func (i *imageResizer) setupOnLoadCb() {
	i.onImgLoadCb = js.FuncOf(func(this js.Value, args []js.Value) any {
		array := args[0]
		i.inBuf = make([]uint8, array.Get("byteLength").Int())
		js.CopyBytesToGo(i.inBuf, array)

		reader := bytes.NewReader(i.inBuf)
		var err error
		var fName string
		i.sourceImg, fName, err = image.Decode(reader)
		if err != nil {
			i.logViaJsConsole(err.Error())
			return nil
		}
		i.logViaJsConsole(fmt.Sprintf("Image loaded, format: %s", fName))
		i.resizeImg()

		return nil
	})
}

func (i *imageResizer) resizeImg() any {
	start := time.Now()
	if i.sourceImg == nil {
		i.logViaJsConsole("Image not loaded")
		return nil
	}

	buf := bytes.NewBuffer(nil)
	thumbnail := imaging.Resize(i.sourceImg, 250, 0, imaging.Lanczos)
	err := jpeg.Encode(buf, thumbnail, &jpeg.Options{Quality: 90})
	if err != nil {
		i.logViaJsConsole(fmt.Sprintf("Cannot create thumbnail: %s", err.Error()))
		return nil
	}

	dst := js.Global().Get("Uint8Array").New(len(buf.Bytes()))
	n := js.CopyBytesToJS(dst, buf.Bytes())
	i.logViaJsConsole(fmt.Sprintf("bytes copied: %d", n))
	js.Global().Call("displayThumbnail", dst)
	i.logViaJsConsole(fmt.Sprintf("time taken: %s", time.Now().Sub(start).String()))
	i.logViaJsConsole("Thumbnail created")
	return nil
}

func (i *imageResizer) logViaJsConsole(msg string) {
	i.console.Call("log", msg)
}
