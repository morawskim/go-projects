package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
)

type contextKey string

const fileNameKey contextKey = "logFileName"

func main() {
	var target string
	counter := atomic.Uint64{}

	flag.StringVar(&target, "target", "", "Target URL for reverse proxy")
	flag.Parse()

	if "" == target {
		log.Fatalln("Target URL cannot be empty")
	}

	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatalln(err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			counter.Add(1)
			fileName := fmt.Sprintf("request_%05d.log", counter.Load())

			reqDump, err := httputil.DumpRequest(request.In, true)
			if err != nil {
				log.Println(err)
			} else {
				err = os.WriteFile(fileName, reqDump, 0644)
				if err != nil {
					log.Println("Error writing request to file:", err)
				}
			}

			ctx := request.Out.Context()
			ctx = context.WithValue(ctx, fileNameKey, fileName)
			request.Out = request.Out.WithContext(ctx)

			request.SetURL(targetUrl)
		},
		ModifyResponse: func(r *http.Response) error {
			fileName, ok := r.Request.Context().Value(fileNameKey).(string)
			if !ok {
				log.Println("Could not find log file name in context")
				return nil
			}

			resDump, err := httputil.DumpResponse(r, true)
			if err != nil {
				return err
			}

			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.WriteString("\n---\n"); err != nil {
				return err
			}
			if _, err := f.Write(resDump); err != nil {
				return err
			}

			return nil
		},
	}

	err = http.ListenAndServe(":8081", proxy)

	if err != nil {
		panic(err)
	}
}
