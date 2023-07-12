package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed frontend/build/*
var content embed.FS

type invokeResponse struct {
	Msg string `json:"msg"`
}

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
	mux.HandleFunc("/invoke", func(w http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		payload, err := io.ReadAll(request.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invokeResponse{Msg: "Unable to read request body"})

			return
		}

		if false == json.Valid(payload) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invokeResponse{Msg: "Payload is not valid JSON}"})

			return
		}

		lambdaInvocationUrl := fmt.Sprintf("%s/2015-03-31/functions/function/invocations", os.Getenv("LAMBDA_RUNTIME_URL"))
		log.Printf("Lambda invoke url: %v", lambdaInvocationUrl)

		newRequest, err := http.NewRequest(
			"POST",
			lambdaInvocationUrl,
			bytes.NewReader(payload),
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invokeResponse{Msg: "Unable to create a request"})

			return
		}

		newRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		response, error := client.Do(newRequest)
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invokeResponse{Msg: "Unable to invoke a function"})

			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)
		w.Write(body)
	})

	fmt.Println("Listing....")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
