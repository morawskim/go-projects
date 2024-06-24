package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
)

//go:embed index.html
var contentFs embed.FS

func main() {
	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(http.FS(contentFs)))
	handler.HandleFunc("/file.json.encrypted", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get(os.Getenv("FILE_URL"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		defer response.Body.Close()

		_, err = io.Copy(w, response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}
}
