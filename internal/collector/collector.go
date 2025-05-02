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

func Start() {
	var c []prometheus.Collector
	var metricsPath = viper.GetString("metrics.path")
	var addr = fmt.Sprintf("%s:%d", viper.GetString("metrics.host"), viper.GetInt("metrics.port"))
	registry := prometheus.NewRegistry()

	targets := viper.GetStringMap("targets")
	for _, config := range targets { // targetIP, config
		if targetConfig, ok := config.(map[string]interface{}); ok {
			collect := targetConfig["collect"].([]interface{})
			for _, item := range collect {
				if metric, ok := item.(map[string]interface{}); ok {
					createCollector(metric, &c)
				}
			}
		}
	}

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	registry.MustRegister(c...)

	// todo health endpoint
	http.Handle(
		metricsPath, promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
				Registry:          registry,
			}),
	)
	logger.Log.Infof("Starting metrics server on %s", addr)
	// To test: curl -H 'Accept: application/openmetrics-text' localhost:8080/metrics
	logger.Log.Fatalln(http.ListenAndServe(addr, nil))
}

func createCollector(m map[string]interface{}, colls *[]prometheus.Collector) {
	switch m["metric_type"] {
	case "gauge":
		rawLabels, ok := m["labels"].([]interface{}) //todo how to make this less ugly?
		if !ok {
			logger.Log.Warnf("Invalid labels format for metric: %s", m["name"])
			return
		}
		labels := make([]string, len(rawLabels))
		for i, label := range rawLabels {
			if strLabel, ok := label.(string); ok {
				labels[i] = strLabel
			} else {
				logger.Log.Warnf("Invalid label type for metric: %s", m["name"])
				return
			}
		}
		g := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: m["name"].(string),
			},
			labels,
		)
		*colls = append(*colls, g)
	default:
		logger.Log.Warnf("Unknown metric type: %s", m["metric_type"])
	}
}

// Collect targets:
// on scrape but
// no more often than given frequency
func Collect() {

}
