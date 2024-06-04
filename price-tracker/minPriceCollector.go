package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type minPriceCollector struct {
	minPrice          *prometheus.Desc
	mapProductToPrice map[string]float64
}

func (m *minPriceCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- m.minPrice
}

func (m *minPriceCollector) Collect(metrics chan<- prometheus.Metric) {
	for product, price := range m.mapProductToPrice {
		metrics <- prometheus.MustNewConstMetric(m.minPrice, prometheus.GaugeValue, price, product)
	}
}

func (m *minPriceCollector) UpdateMinPrice(product string, price float64) {
	if _, ok := m.mapProductToPrice[product]; !ok {
		m.mapProductToPrice[product] = price
	}

	if val, ok := m.mapProductToPrice[product]; ok {
		if val > price {
			m.mapProductToPrice[product] = price
		}
	}
}

func (m *minPriceCollector) GetMinPrice(product string) (float64, bool) {
	price, ok := m.mapProductToPrice[product]

	return price, ok
}

func newMinPriceCollector() *minPriceCollector {
	return &minPriceCollector{
		mapProductToPrice: make(map[string]float64),
		minPrice: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "min_price"),
			"Minimal price which has been seen",
			[]string{"product"},
			nil,
		),
	}
}
