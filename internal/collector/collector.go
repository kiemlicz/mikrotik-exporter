package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"mikrotik-exporter/internal/logger"
	"net/http"
)

var (
	log = logger.Log
)

func Start() {
	var c = []prometheus.Collector{}
	var metricsPath = viper.GetString("metrics.path")
	var addr = fmt.Sprintf("%s:%d", viper.GetString("metrics.host"), viper.GetInt("metrics.port"))
	var metrics = viper.Get("metrics.collect").([]interface{})
	registry := prometheus.NewRegistry()
	for _, e := range metrics {
		m := e.(map[string]interface{})
		switch m["metric_type"] {
		case "gauge":
			log.Info("gauge")
		default:
			log.Warnf("Unknown metric type: %s", m["metric_type"])
		}
	}

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	registry.MustRegister(c...)

	http.Handle(
		metricsPath, promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
				Registry:          registry,
			}),
	)
	// To test: curl -H 'Accept: application/openmetrics-text' localhost:8080/metrics
	log.Fatalln(http.ListenAndServe(addr, nil))
}
