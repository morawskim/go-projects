package main

import (
	_ "embed"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"html/template"
	"log/slog"
	"net/http"
	"time"
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

//go:embed index.html
var tmpl string

func updateLastScrapeMetric() {
	lastScrape.SetToCurrentTime()
}

func registerMetrics(items []item2, minPriceCollector *minPriceCollector) {
	prometheus.MustRegister(priceMetric)
	prometheus.MustRegister(productScraper)
	prometheus.MustRegister(lastScrape)
	for _, i := range items {
		productScraper.With(prometheus.Labels{"Product": i.productName}).Set(0)
	}
	prometheus.MustRegister(minPriceCollector)
}

func createChannel(minPriceCollector *minPriceCollector) chan metric {
	ch := make(chan metric)

	go func() {
		for m := range ch {
			priceMetric.With(prometheus.Labels{"Product": m.product}).Set(m.price)
			productScraper.With(prometheus.Labels{"Product": m.product}).Set(1)
			minPriceCollector.UpdateMinPrice(m.product, m.price)
		}
	}()

	return ch
}

type StatusItem struct {
	Product string
	Url     string
	Price   float64
	Status  int
}

func register(products []item2) {
	t := template.Must(template.New("status").Parse(tmpl))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		productStatus := make([]StatusItem, 0, len(products))

		lastScrapeMetricDto := io_prometheus_client.Metric{}
		lastScrape.Write(&lastScrapeMetricDto)

		for _, p := range products {
			priceMetricValue, err := extractGaugeMetricValue(p.productName, priceMetric)
			if err != nil {
				slog.Default().Error(
					fmt.Sprintf("unable to get price metric for product %v. reason: %s", p.productName, err.Error()),
					slog.String("product", p.productName),
				)
				productStatus = append(productStatus, createStatusItemWithError(p.productName, p.productUrl))
				continue
			}

			scrapeMetricValue, err := extractGaugeMetricValue(p.productName, productScraper)
			if err != nil {
				slog.Default().Error(
					fmt.Sprintf("unable to get product scraper metric for product %v. reason: %s", p.productName, err.Error()),
					slog.String("product", p.productName),
				)
				productStatus = append(productStatus, createStatusItemWithError(p.productName, p.productUrl))
				continue
			}

			productStatus = append(productStatus, StatusItem{
				Product: p.productName,
				Url:     p.productUrl,
				Price:   priceMetricValue,
				Status:  int(scrapeMetricValue),
			})
		}

		err := t.Execute(w, struct {
			LastScrape string
			Products   []StatusItem
		}{
			LastScrape: time.Unix(int64(*lastScrapeMetricDto.Gauge.Value), 0).Format("2006-01-02 15:04:05 -07:00"),
			Products:   productStatus,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func extractGaugeMetricValue(productName string, metric *prometheus.GaugeVec) (float64, error) {
	dto := io_prometheus_client.Metric{}
	g, e := metric.GetMetricWithLabelValues(productName)
	if e != nil {
		return 0, e
	}
	g.Write(&dto)

	return *dto.Gauge.Value, nil
}

func createStatusItemWithError(productName, productUrl string) StatusItem {
	return StatusItem{
		Product: productName,
		Url:     productUrl,
		Price:   0,
		Status:  0,
	}
}
