package action

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"net/http"
)

type EncryptPDFRequest struct {
	File    string
	UserPw  string `json:"user_pw"`
	OwnerPw string `json:"owner_pw"`
}

type EncryptPDFResponse struct {
	File string `json:"file"`
}

func HandleEncryptPDF(w http.ResponseWriter, req *http.Request) {
	var err error
	action := &EncryptPDFRequest{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(action)

	if err != nil {
		http.Error(w, "Cannot decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(action.File) == 0 {
		http.Error(w, "The file field cannot be empty", http.StatusBadRequest)
		return
	}

	if len(action.OwnerPw) == 0 || len(action.UserPw) == 0 {
		http.Error(w, "The password cannot be empty", http.StatusBadRequest)
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(action.File)
	if err != nil {
		http.Error(w, "Cannot decode file: "+err.Error(), http.StatusBadRequest)
		return
	}

	if mtype := mimetype.Detect(rawDecodedText); !mtype.Is("application/pdf") {
		http.Error(w, fmt.Sprintf("Expect PDF file got %q", mtype.String()), http.StatusBadRequest)
		return
	}

	out := new(bytes.Buffer)
	conf := pdfcpu.NewAESConfiguration(action.UserPw, action.OwnerPw, 256)
	err = api.Encrypt(bytes.NewReader(rawDecodedText), out, conf)

	if err != nil {
		http.Error(w, "Cannot encrypt PDF file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := EncryptPDFResponse{
		File: base64.StdEncoding.EncodeToString(out.Bytes()),
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
