package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed frontend/build/*
var content embed.FS

func main() {

	//sdkConfig, err := config.LoadDefaultConfig(context.Background())
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load AWS default config: %s", err))
	//}
	//s3Client := s3.NewFromConfig(sdkConfig)
	//
	//s3Client.ListO

	fsys := fs.FS(content)
	html, _ := fs.Sub(fsys, "frontend/build")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(html)))

	fmt.Println("Listing....")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
