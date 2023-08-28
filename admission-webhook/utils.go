package main

import (
	"encoding/json"
	"fmt"
	"io"
	admission "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"net/http"
)

func sendResponse(w http.ResponseWriter, obj *admission.AdmissionReview) {
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

func extract(input io.ReadCloser) (*admission.AdmissionReview, error) {
	body, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %w", err)
	}
	defer input.Close()

	admissionReview := admission.AdmissionReview{}
	if err = json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("cannot unmarshal AdmissionReview: %w", err)
	}

	return &admissionReview, nil
}

func isPodResource(admissionReview *admission.AdmissionReview) bool {
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admissionReview.Request.Resource != podResource {
		return false
	}

	return true
}

func createAdmissionReview(response *admission.AdmissionResponse) *admission.AdmissionReview {
	responseAdmissionReview := &admission.AdmissionReview{}
	responseAdmissionReview.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "admission.k8s.io",
		Version: "v1",
		Kind:    "AdmissionReview",
	})
	responseAdmissionReview.Response = response

	return responseAdmissionReview
}
