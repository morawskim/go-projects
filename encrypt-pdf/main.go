package main

import (
	"github.com/gorilla/mux"
	"github.com/morawskim/go-projects/encrypt-pdf/action"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/encrypt", action.HandleEncryptPDF).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()
}
