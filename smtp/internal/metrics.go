package internal

import (
	"github.com/prometheus/client_golang/prometheus"
)

var buildInfoMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "gosmtp_build_info",
	Help: "Build information",
}, []string{"version", "commit", "date"})

var forwardedMails = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "gosmtp_forwarded_mails",
	Help: "Number of forwarded mails",
}, []string{"status"})

func MustRegisterMetrics() {
	prometheus.MustRegister(buildInfoMetric)
	prometheus.MustRegister(forwardedMails)
	buildInfoMetric.WithLabelValues(Version, Commit, BuildDate).Set(1)
}

func IncreaseForwardedMailsMetric(hasError bool) {
	if hasError {
		forwardedMails.WithLabelValues("error").Inc()
	} else {
		forwardedMails.WithLabelValues("success").Inc()
	}
}
