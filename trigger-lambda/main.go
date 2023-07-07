package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

//go:embed index.html
var indexPage []byte

func main() {

	//sdkConfig, err := config.LoadDefaultConfig(context.Background())
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load AWS default config: %s", err))
	//}
	//s3Client := s3.NewFromConfig(sdkConfig)
	//
	//s3Client.ListO

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(indexPage)
	})
	fmt.Println("Listing....")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
