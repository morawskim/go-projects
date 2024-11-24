package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type upstreamResponse struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

var requestTimeout = 100 * time.Second
var tickerInterval = 10 * time.Second

func main() {
	httpClient := &http.Client{Timeout: requestTimeout + 10*time.Second}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flusher, isFlusher := w.(http.Flusher)
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		upstreamUrl, err := url.Parse(r.URL.Query().Get("url"))
		if err != nil {
			err := sendJsonResponse[errorResponse](
				w,
				http.StatusBadRequest,
				buildErrorResponse(fmt.Errorf("invalid URL: %v", err)),
			)
			if err != nil {
				slog.Error("unable to send error response: " + err.Error())
			}
			return
		}

		headers := r.Header.Clone()
		headers.Del("Host")
		headers.Add("X-Forwarded-For", r.RemoteAddr)
		proxyRequest, err := http.NewRequestWithContext(ctx, r.Method, upstreamUrl.String(), io.NopCloser(strings.NewReader("")))

		if err != nil {
			slog.Default().Info("unable to create proxy request: " + err.Error())
			err := sendJsonResponse[errorResponse](
				w,
				http.StatusInternalServerError,
				buildErrorResponse(fmt.Errorf("unable to create request")),
			)
			if err != nil {
				slog.Error("unable to send error response: " + err.Error())
			}
			return
		}

		proxyRequest.Header = headers
		ch := sendRequest(httpClient, proxyRequest)

		timer := time.NewTicker(tickerInterval)
		defer timer.Stop()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		for {
			select {
			case fResponse := <-ch:
				err := encodeJson(w, fResponse)
				if err != nil {
					slog.Default().Info("unable to send response: " + err.Error())
				}
				return
			case <-ctx.Done():
				err := encodeJson(w, ctx.Err())
				if err != nil {
					slog.Default().Info("unable to send response: " + err.Error())
				}
				return
			case <-timer.C:
				_, err := w.Write([]byte("\n"))
				if err != nil {
					slog.Default().Info("unable to send response: " + err.Error())
				}

				if err == nil && isFlusher {
					flusher.Flush()
				}
			}
		}
	})

	http.ListenAndServe("localhost:8080", http.DefaultServeMux)
}

func buildErrorResponse(err error) errorResponse {
	return errorResponse{
		Error: err.Error(),
	}
}

func encodeJson[T any](w http.ResponseWriter, body T) error {
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return fmt.Errorf("cannot encode response: %v", err)
	}

	return nil
}

func sendJsonResponse[T any](w http.ResponseWriter, httpStatusCode int, body T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	return encodeJson(w, body)
}

func sendRequest(httpClient *http.Client, proxyRequest *http.Request) chan upstreamResponse {
	ch := make(chan upstreamResponse)
	go func() {
		do, err := httpClient.Do(proxyRequest)
		if err != nil {
			panic(fmt.Errorf("http client error: %v", err))
		}
		defer do.Body.Close()

		body, err := io.ReadAll(do.Body)
		if err != nil {
			panic(fmt.Errorf("cannot read body: %v", err))
		}

		fw := upstreamResponse{
			StatusCode: do.StatusCode,
			Headers:    do.Header,
			Body:       string(body),
		}

		ch <- fw
		close(ch)
	}()

	return ch
}
