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

var priceMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "price_tracker",
	Help: "Trace price of product",
}, []string{"Product"})

var productScraper = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "product_scraper",
	Help: "Status of the scraping",
}, []string{"Product"})

var lastScrape = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "last_scrape_time_seconds",
	Help: "Last scrape of prices",
})

func updateLastScrapeMetric() {
	lastScrape.SetToCurrentTime()
}

func registerMetrics(items []item2) {
	prometheus.MustRegister(priceMetric)
	prometheus.MustRegister(productScraper)
	prometheus.MustRegister(lastScrape)
	for _, i := range items {
		productScraper.With(prometheus.Labels{"Product": i.productName}).Set(0)
	}
}

func createChannel() chan metric {
	ch := make(chan metric)

	go func() {
		for m := range ch {
			priceMetric.With(prometheus.Labels{"Product": m.product}).Set(m.price)
			productScraper.With(prometheus.Labels{"Product": m.product}).Set(1)
		}
	}()

	return ch
}

func register() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
