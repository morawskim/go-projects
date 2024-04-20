package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type metric struct {
	product string
	price   float64
}

func registerMetrics(items []item2) {
	prometheus.MustRegister(visitedURLs)
	for _, i := range items {
		visitedURLs.With(prometheus.Labels{"Product": i.productName}).Set(-1)
	}
}

func createChannel() chan metric {
	ch := make(chan metric)

	go func() {
		for m := range ch {
			visitedURLs.With(prometheus.Labels{"Product": m.product}).Set(m.price)
		}
	}()

	return ch
}

func register() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
