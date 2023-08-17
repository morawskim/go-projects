package main

import (
	"encoding/json"
	"flag"
	"io"
	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"net/http"
)

func main() {
	var tlsKey, tlsCert string
	flag.StringVar(&tlsKey, "tlsKey", "/etc/certs/tls.key", "Path to the TLS key")
	flag.StringVar(&tlsCert, "tlsCert", "/etc/certs/tls.crt", "Path to the TLS certificate")
	flag.Parse()

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Validating")

		all, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln("Cannot read request body", err)
			return
		}
		defer r.Body.Close()

		log.Println("request body", string(all))

		admissionReview := admission.AdmissionReview{}
		err = json.Unmarshal(all, &admissionReview)
		if err != nil {
			log.Println("unmarshal error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseObj := admission.AdmissionResponse{Allowed: true, UID: admissionReview.Request.UID}

		responseAdmissionReview := &admission.AdmissionReview{}
		responseAdmissionReview.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "admission.k8s.io",
			Version: "v1",
			Kind:    "AdmissionReview",
		})
		responseAdmissionReview.Response = &responseObj

		log.Println("sending response", responseObj.Allowed, responseObj.UID)
		respBytes, err := json.Marshal(responseAdmissionReview)

		if err != nil {
			log.Println("response marshal error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(respBytes); err != nil {
			log.Println("response write error", err)
		}
	})

	log.Println("Starting server on port 4443")
	err := http.ListenAndServeTLS(":4443", tlsCert, tlsKey, nil)
	if err != nil {
		log.Panicln(err)
	}
}
