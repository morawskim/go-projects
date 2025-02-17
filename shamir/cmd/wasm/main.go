package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/shamir"
	"strings"
	"syscall/js"
)

const SHAMIR_PARTS = 5
const SHARIM_THRESHOLD = 3

func main() {
	js.Global().Set("shamirSplit", wrapperForShamirSplit())
	js.Global().Set("shamirCombine", wrapperForShamirCombine())
	//otherwise we get error "Go program has already exited" in web browser
	select {}
}

func wrapperForShamirSplit() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		result := map[string]any{"error": nil, "result": nil}

		if len(args) != 1 {
			result["error"] = "Invalid number of arguments passed"
			return result
		}

		value := args[0].String()
		parts, err := split([]byte(value))

		if err == nil {
			jsonData, err := json.Marshal(parts)

			if err == nil {
				result["result"] = string(jsonData)
				result["error"] = nil
			} else {
				result["result"] = nil
				result["error"] = fmt.Errorf("failed to encode shamir parts: %s", err)
			}
		} else {
			result["result"] = nil
			result["error"] = err
		}

		return result
	})
}

func wrapperForShamirCombine() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		result := map[string]any{"error": nil, "result": nil}

		if len(args) != 1 {
			result["error"] = "Invalid number of arguments passed"
			return result
		}

		value := args[0].String()
		decodedString, err := combine(strings.Split(value, "\n"))

		if err == nil {
			result["result"] = decodedString
			result["error"] = nil
		} else {
			result["result"] = nil
			result["error"] = err
		}

		return result
	})
}

func split(secretBuf []byte) ([]string, error) {
	byteParts, err := shamir.Split(secretBuf, SHAMIR_PARTS, SHARIM_THRESHOLD)
	hexParts := make([]string, 0, len(byteParts))

	if err != nil {
		return nil, err
	}

	for _, bytePart := range byteParts {
		hexParts = append(hexParts, fmt.Sprintf("%x", bytePart))
	}

	return hexParts, nil
}

func combine(hexParts []string) (string, error) {
	var byteParts [][]byte
	for _, hexPart := range hexParts {
		b, err := hex.DecodeString(hexPart)
		if err != nil {
			return "", fmt.Errorf("Failed to decode %q: %v\n", hexPart, err)
		}
		byteParts = append(byteParts, b)
	}
	secretBytes, err := shamir.Combine(byteParts)
	if err != nil {
		return "", fmt.Errorf("Failed to combine secret: %v\n", err)
	}

	return string(secretBytes), nil
}
