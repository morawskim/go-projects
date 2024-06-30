package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"syscall/js"
)

func main() {
	fmt.Println("Hello World from Webassembly")
	imgResizer := newImageResizer()
	imgResizer.setupOnLoadCb()

	js.Global().Set("base64Encode", wrapperForBase64())
	js.Global().Set("loadImage", imgResizer.onImgLoadCb)

	// see https://github.com/golang/go/issues/41310
	js.Global().Set("goFetch", js.FuncOf(func(this js.Value, args []js.Value) any {
		url := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				r, err := http.Get(url)
				if err != nil {
					reject.Invoke(fmt.Sprintf("cannot send request: %s", err.Error()))
					return
				}

				defer r.Body.Close()
				body, err := io.ReadAll(r.Body)

				if err != nil {
					reject.Invoke(fmt.Sprintf("cannot read body: %s", err.Error()))
					return
				}

				resolve.Invoke(string(body))
			}()

			return nil
		})
		// Create and return the Promise
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	}))

	//otherwise we get error "Go program has already exited" in web browser
	select {}
}

func encodeStringToBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func wrapperForBase64() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		result := map[string]any{"error": nil, "result": nil}

		if len(args) != 1 {
			result["error"] = "Invalid no of arguments passed"
			return result
		}

		value := args[0].String()
		base64String := encodeStringToBase64(value)
		result["result"] = base64String

		return result
	})
}
