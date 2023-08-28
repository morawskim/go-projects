package main

import (
	"encoding/json"
	admission "k8s.io/api/admission/v1"
	"log"
	"net/http"
)

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Mutate request")

	admissionReview, err := extract(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot extract AdmissionReview from request", http.StatusBadRequest)
		return
	}

	if !isPodResource(admissionReview) {
		log.Printf("expect resource to be pod, got %s\n", admissionReview.Request.Resource)
		return
	}

	pt := admission.PatchTypeJSONPatch
	patch := `[{ "op": "add", "path": "/metadata/labels/foo", "value": "bar" }]`

	admissionResponse := admission.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		PatchType: &pt,
		Patch:     []byte(patch),
	}
	obj := createAdmissionReview(&admissionResponse)

	respBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println("response marshal error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		log.Println("response write error", err)
	}
}
