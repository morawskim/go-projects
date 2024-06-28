package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello World from Webassembly")
	js.Global().Set("base64Encode", wrapperForBase64())

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
