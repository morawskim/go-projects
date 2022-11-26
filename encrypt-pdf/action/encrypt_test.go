package action

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandleEmptyBody(t *testing.T) {
	buf := new(bytes.Buffer)
	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	// We should get a good status code
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "Cannot decode JSON: EOF", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleInvalidJSONBody(t *testing.T) {
	buf := strings.NewReader("{invalid")

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	// We should get a good status code
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "Cannot decode JSON: invalid character", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleEmptyFileField(t *testing.T) {
	buf := new(bytes.Buffer)
	actionRequest := EncryptPDFRequest{
		File: "",
	}
	body, _ := json.Marshal(actionRequest)
	buf.Write(body)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "The file field cannot be empty", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleInvalidBase64(t *testing.T) {
	buf := new(bytes.Buffer)
	actionRequest := EncryptPDFRequest{
		File:    "a",
		UserPw:  "foo",
		OwnerPw: "bar",
	}
	body, _ := json.Marshal(actionRequest)
	buf.Write(body)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "Cannot decode file: illegal base64 data", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleInvalidPDFFile(t *testing.T) {
	buf := new(bytes.Buffer)
	actionRequest := EncryptPDFRequest{
		File:    "Zm9v",
		UserPw:  "foo",
		OwnerPw: "bar",
	}
	body, _ := json.Marshal(actionRequest)
	buf.Write(body)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "Expect PDF file got", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleEmptyPassword(t *testing.T) {
	buf := new(bytes.Buffer)
	actionRequest := EncryptPDFRequest{
		File:    "Zm9v",
		UserPw:  "",
		OwnerPw: "",
	}
	body, _ := json.Marshal(actionRequest)
	buf.Write(body)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	if want, got := "The password cannot be empty", w.Body.String(); !strings.Contains(got, want) {
		t.Fatalf("expected a %q, instead got: %q", want, got)
	}
}

func TestHandleEncryptPDF(t *testing.T) {
	file, _ := os.ReadFile("../fixture/test.pdf")
	buf := new(bytes.Buffer)
	actionRequest := EncryptPDFRequest{
		File:    base64.StdEncoding.EncodeToString(file),
		UserPw:  "foo",
		OwnerPw: "bar",
	}
	body, _ := json.Marshal(actionRequest)
	buf.Write(body)

	req := httptest.NewRequest(http.MethodPost, "http://example.com", buf)
	w := httptest.NewRecorder()

	HandleEncryptPDF(w, req)

	// We should get a good status code
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d (%v)", want, got, w.Body.String())
	}
}
