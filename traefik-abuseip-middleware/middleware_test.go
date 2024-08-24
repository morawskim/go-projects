package traefik_abuseip_middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		remoteAddr         string
		expectedStatusCode int
	}{
		{name: "ip allowed", expectedStatusCode: http.StatusOK, remoteAddr: "1.1.5.5:54555"},
		{name: "ip forbidden", expectedStatusCode: http.StatusForbidden, remoteAddr: "1.1.2.2:54555"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := httptest.NewServer(http.FileServer(http.Dir("./test-data")))
			defer srv.Close()

			cfg := CreateConfig()
			cfg.AbuseIpFile = srv.URL + "/blocked-ips.txt"

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := New(ctx, next, cfg, "demo")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			req.RemoteAddr = test.remoteAddr
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			if recorder.Code != test.expectedStatusCode {
				t.Fatalf("got %d, want %d", recorder.Code, test.expectedStatusCode)
			}
		})
	}
}
