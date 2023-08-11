package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
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

var basePath = ""

func main() {

	//sdkConfig, err := config.LoadDefaultConfig(context.Background())
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load AWS default config: %s", err))
	//}
	//s3Client := s3.NewFromConfig(sdkConfig)
	//
	//s3Client.ListO

	flag.StringVar(&basePath, "base-path", "", "The URL prefix for all requests, relative to the host root.")
	flag.Parse()

	fsys := fs.FS(content)
	html, _ := fs.Sub(fsys, "frontend/build")

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix(basePath, &myFsHandler{
		internalHandler: http.FileServer(http.FS(html)),
		template:        template.Must(template.ParseFS(html, "index.html")),
		basePath:        basePath,
	}))
	mux.HandleFunc(basePath+"/invoke", func(w http.ResponseWriter, request *http.Request) {
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
		response, err := client.Do(newRequest)
		if err != nil {
			log.Println("Invoke Lambda error", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invokeResponse{Msg: "Unable to invoke a function"})

			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)
		w.Write(body)
	})

	fmt.Println("Listing on port :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
