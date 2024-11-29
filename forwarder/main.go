package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

type upstreamResponse struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
}

type sendRequestResult struct {
	err    error
	result upstreamResponse
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

	httpAddr := flag.String("addr", "localhost:8080", "http listen address")
	flag.Parse()

	server := &http.Server{
		Addr:    *httpAddr,
		Handler: http.DefaultServeMux,
	}

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

		if validErr := validateURL(upstreamUrl.String()); validErr != nil {
			err := sendJsonResponse[errorResponse](
				w,
				http.StatusBadRequest,
				buildErrorResponse(validErr),
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
				if fResponse.err != nil {
					slog.Info("unable to send request: " + fResponse.err.Error())

					var urlErr *url.Error
					if errors.As(fResponse.err, &urlErr) {
						if urlErr.Timeout() {
							err := encodeJson(w, buildErrorResponse(fmt.Errorf("timeout")))
							if err != nil {
								slog.Default().Info("unable to send response: " + err.Error())
							}
							return
						}
					}

					err := encodeJson(w, buildErrorResponse(fmt.Errorf("send request failed")))
					if err != nil {
						slog.Default().Info("unable to send response: " + err.Error())
					}
				} else {
					err := encodeJson(w, fResponse.result)
					if err != nil {
						slog.Default().Info("unable to send response: " + err.Error())
					}
				}
				return
			case <-ctx.Done():
				err := encodeJson(w, buildErrorResponse(errors.New("timeout")))
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		slog.Info(fmt.Sprintf("Starting server on %s", *httpAddr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error: %v", err)
		}
	}()

	// Block until the context is canceled (signal received)
	<-ctx.Done()
	slog.Info("Shutting down server...")

	// Create a timeout context for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Warn("Server forced to shut down: %v", err)
	} else {
		slog.Info("Server stopped gracefully.")
	}
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

func sendRequest(httpClient *http.Client, proxyRequest *http.Request) chan sendRequestResult {
	ch := make(chan sendRequestResult)
	go func() {
		defer close(ch)
		do, err := httpClient.Do(proxyRequest)
		if err != nil {
			ch <- sendRequestResult{err: err, result: upstreamResponse{}}
			return
		}
		defer do.Body.Close()

		body, err := io.ReadAll(do.Body)
		if err != nil {
			ch <- sendRequestResult{
				err:    fmt.Errorf("cannot read body: %v", err),
				result: upstreamResponse{},
			}
			return
		}

		ch <- sendRequestResult{
			err: nil,
			result: upstreamResponse{
				StatusCode: do.StatusCode,
				Headers:    do.Header,
				Body:       string(body),
			},
		}
	}()

	return ch
}

func validateURL(input string) error {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("URL must have a scheme (e.g., http, https) and a host")
	}

	return nil
}
