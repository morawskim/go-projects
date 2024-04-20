package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/richardjennings/tapo/pkg/tapo"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type TapoMetrics struct {
	CurrentPower float64
	TodayEnergy  float64
}

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	var interval time.Duration
	flag.DurationVar(&interval, "interval", 60*time.Second, "")
	flag.Parse()

	metricCurrentPower := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "app",
		Name:      "current_power",
		Help:      "Current power usage ",
	})

	metricPowerEnergy := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "app",
		Name:      "power_energy",
		Help:      "Today energy usage in Wh",
	})

	reg := prometheus.NewRegistry()
	reg.MustRegister(metricCurrentPower, metricPowerEnergy)

	t, err := tapo.NewTapo(os.Getenv("TAPO_IP"), os.Getenv("TAPO_USERNAME"), os.Getenv("TAPO_PASSWORD"))
	if err != nil {
		panic(err)
	}

	go func() {
		// configure this value via command options default 60s
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				logger.Debug("Time to fetch metrics from tapo")
				updateMetrics(metricCurrentPower, metricPowerEnergy, t)
			}
		}
	}()

	updateMetrics(metricCurrentPower, metricPowerEnergy, t)
	err = http.ListenAndServe(":8080", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	if err != nil {
		panic(err)
	}
}

func fetchDataFromTapo(tapo *tapo.Tapo) (TapoMetrics, error) {
	res, err := tapo.GetEnergyUsage()

	if err != nil {
		return TapoMetrics{}, err
	}

	if res["error_code"] != float64(0) {
		return TapoMetrics{}, errors.New("non zero error code from taapo")
	}

	//fmt.Println(res["result"].(map[string]interface{})["today_runtime"].(float64))
	//fmt.Println(res["result"].(map[string]interface{})["local_time"].(string))

	return TapoMetrics{
		CurrentPower: res["result"].(map[string]interface{})["current_power"].(float64),
		TodayEnergy:  res["result"].(map[string]interface{})["today_energy"].(float64),
	}, nil
}

func updateMetrics(metricCurrentPower, metricPowerEnergy prometheus.Gauge, t *tapo.Tapo) {
	metrics, err := fetchDataFromTapo(t)

	if err != nil {
		logger.Error(fmt.Sprintf("fetching metrics from tapo fail - %q", err))
	}

	metricCurrentPower.Set(metrics.CurrentPower)
	metricPowerEnergy.Set(metrics.TodayEnergy)
}
