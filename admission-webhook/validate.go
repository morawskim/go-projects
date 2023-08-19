package main

import (
	"encoding/json"
	"io"
	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Validating request")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Cannot read request body", err)
		http.Error(w, "Cannot read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//log.Println("request body", string(body))
	admissionReview := admission.AdmissionReview{}
	if err = json.Unmarshal(body, &admissionReview); err != nil {
		log.Println("unmarshal request error", err)
		http.Error(w, "Decode request body failed", http.StatusInternalServerError)
		return
	}

	deploymentResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admissionReview.Request.Resource != deploymentResource {
		log.Printf("expect resource to be %s\n", deploymentResource)
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
}
