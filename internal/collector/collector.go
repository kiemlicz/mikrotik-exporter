package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	customMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "custom_metric_name",
			Help: "Description of the custom metric",
		},
		[]string{"label1", "label2"},
	)
)

func init() {
	// Register the custom metric with Prometheus
	prometheus.MustRegister(customMetric)
}

func StartMetricsServer() {
	// Expose metrics on /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
