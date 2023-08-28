package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var tlsKey, tlsCert string
	flag.StringVar(&tlsKey, "tlsKey", "/etc/certs/tls.key", "Path to the TLS key")
	flag.StringVar(&tlsCert, "tlsCert", "/etc/certs/tls.crt", "Path to the TLS certificate")
	flag.Parse()

	http.HandleFunc("/validate", validateHandler)
	http.HandleFunc("/mutate", mutateHandler)

	log.Println("Starting server on port 4443")
	err := http.ListenAndServeTLS(":4443", tlsCert, tlsKey, nil)
	if err != nil {
		log.Panicln(err)
	}
}
