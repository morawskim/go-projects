package main

import (
	"encoding/json"
	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Validating request")

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

	raw := admissionReview.Request.Object.Raw
	pod := corev1.Pod{}

	if err := json.Unmarshal(raw, &pod); err != nil {
		log.Println("Decode pod JSON failed", err)
		http.Error(w, "Decode pod JSON failed", http.StatusInternalServerError)
	}

	var responseObj admission.AdmissionResponse
	_, exists := pod.GetLabels()["ContactPerson"]
	if exists {
		responseObj = admission.AdmissionResponse{Allowed: true, UID: admissionReview.Request.UID}
	} else {
		responseObj = admission.AdmissionResponse{
			Allowed: false,
			UID:     admissionReview.Request.UID,
			Result:  &metav1.Status{Message: "Label ContactPerson not found"},
		}
	}

	obj := createAdmissionReview(&responseObj)
	sendResponse(w, obj)
}
