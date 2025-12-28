package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const TargetUrlQueryParam string = "__targetUrl"

func main() {
	port := flag.Uint("port", 8080, "Port to listen on")
	delay := flag.Duration("delay", 10*time.Second, "Delay before forwarding request")

	flag.Parse()

	proxy := &httputil.ReverseProxy{
		Rewrite:  rewriteRequestBuilder(*delay),
		Director: nil,
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !r.URL.Query().Has(TargetUrlQueryParam) {
				http.Error(
					w,
					"Missing required query parameter: __targetUrl",
					http.StatusBadRequest,
				)
				return
			}

			targetUrlParam := r.URL.Query().Get(TargetUrlQueryParam)

			if targetUrlParam == "" {
				http.Error(
					w,
					"__targetUrl query parameter cannot be empty",
					http.StatusBadRequest,
				)
				return
			}

			_, err := url.Parse(targetUrlParam)
			if err != nil {
				http.Error(
					w,
					"__targetUrl query parameter is not valid URL",
					http.StatusBadRequest,
				)
				return
			}

			proxy.ServeHTTP(w, r)
		}),
	}

	log.Printf("Listening on port %d\n", *port)
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

func rewriteRequestBuilder(delay time.Duration) func(*httputil.ProxyRequest) {
	return func(pr *httputil.ProxyRequest) {
		time.Sleep(delay)

		if !pr.Out.URL.Query().Has(TargetUrlQueryParam) {
			return
		}

		targetParam := pr.In.URL.Query().Get(TargetUrlQueryParam)
		target, err := url.Parse(targetParam)

		if err != nil {
			return
		}

		targetQuery := target.RawQuery
		pr.Out.URL.Scheme = target.Scheme
		pr.Out.URL.Host = target.Host
		pr.Out.URL.Path, pr.Out.URL.RawPath = target.Path, target.RawPath
		pr.Out.Host = pr.Out.URL.Host

		if targetQuery == "" || pr.Out.URL.RawQuery == "" {
			pr.Out.URL.RawQuery = targetQuery + pr.Out.URL.RawQuery
		} else {
			pr.Out.URL.RawQuery = targetQuery + "&" + pr.Out.URL.RawQuery
		}
		return
	}
}
