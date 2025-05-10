package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"mikrotik-exporter/internal/logger"
	"net/http"
	"sync"
)

type MikrotikMetric struct {
	desc       *prometheus.Desc
	metricType prometheus.ValueType
	request    string //todo different type?
	parse      string
}

type MikrotikDeviceCollector struct {
	device     *MikrotikDevice
	deviceData sync.Map
	metrics    map[string]MikrotikMetric
}

func NewMikrotikDeviceCollector(host string, port int, username, password string, collect []interface{}) MikrotikDeviceCollector {
	m := make(map[string]MikrotikMetric)
	for _, item := range collect {
		if metric, ok := item.(map[string]interface{}); ok {
			name := metric["name"].(string)
			labels := metric["labels"].([]string)
			desc := prometheus.NewDesc(
				name,
				"",
				labels,
				nil,
			)
			valueType, err := parseValueType(metric["metric_type"].(string))
			if err != nil {
				// fixme panic?
				return MikrotikDeviceCollector{}
			}
			request := metric["request"].(string)
			parse := metric["parse"].(string)
			m[name] = MikrotikMetric{
				desc:       desc,
				metricType: valueType,
				request:    request,
				parse:      parse,
			}
		}
	}
	device := &MikrotikDevice{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	return MikrotikDeviceCollector{
		device:     device,
		deviceData: sync.Map{},
		metrics:    m,
	}
}

func (c *MikrotikDeviceCollector) Collect() { //same signature as proms?
	//
}

type MikrotikCollector struct {
	devices []MikrotikDeviceCollector
}

func (collector *MikrotikCollector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, descs)
}

// Collect called every scrape interval
func (collector *MikrotikCollector) Collect(ch chan<- prometheus.Metric) {

	//will force refresh if stale data

	for _, d := range collector.devices {

	}
}

func Start() {
	d := make([]MikrotikDeviceCollector, 0)
	metricsPath := viper.GetString("metrics.path")
	addr := fmt.Sprintf("%s:%d", viper.GetString("metrics.host"), viper.GetInt("metrics.port"))
	registry := prometheus.NewPedanticRegistry()
	targets := viper.GetStringMap("targets")

	for targetIP, config := range targets { // targetIP, config
		if targetConfig, ok := config.(map[string]interface{}); ok {
			collect := targetConfig["collect"].([]interface{})
			d = append(d, NewMikrotikDeviceCollector( //isn't it copied?
				targetIP,
				viper.GetInt("targets.port"),
				viper.GetString("targets.username"),
				viper.GetString("targets.password"),
				collect,
			))
		}
	}

	mikrotikCollector := MikrotikCollector{
		devices: d,
	}

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	registry.MustRegister(&mikrotikCollector)

	// todo health endpoint
	http.Handle(
		metricsPath, promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)
	logger.Log.Infof("Starting metrics server on %s", addr)
	// To test: curl -H 'Accept: application/openmetrics-text' localhost:8080/metrics
	logger.Log.Fatalln(http.ListenAndServe(addr, nil))
}
