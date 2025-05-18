package collector

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"mikrotik-exporter/internal/logger"
	"net/http"
	"time"
)

type MikrotikMetric struct {
	desc       *prometheus.Desc
	metricType prometheus.ValueType

	lastRefresh time.Time
	lastValue   map[string]string
}

func NewMikrotikMetric(
	desc *prometheus.Desc,
	metricType prometheus.ValueType,
	connector *MikrotikConnector,
	request MetricRequest,
	ctx context.Context,
	interval time.Duration,
) *MikrotikMetric {
	m := &MikrotikMetric{
		desc:       desc,
		metricType: metricType,
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//lock?
				m.Refresh(connector, request)
			case <-ctx.Done():
				return
			}
		}
	}()
	return m
}

func (m *MikrotikMetric) Refresh(connector *MikrotikConnector, request MetricRequest) {
	connector.Request(request)

	if err == nil {
		// parse the response
		m.lastValue = string(response)
		m.lastRefresh = time.Now()
	} else {
		logger.Log.Errorf("Failed to fetch data: %v", err)
	}

}

type MikrotikDeviceCollector struct {
	connector   *MikrotikConnector
	minInterval time.Duration
	metrics     map[string]*MikrotikMetric
}

func NewMikrotikDeviceCollector(
	host string,
	port int,
	username string,
	password string,
	useTls bool,
	interval time.Duration,
	collect []interface{},
) MikrotikDeviceCollector {
	m := make(map[string]*MikrotikMetric)

	connector := NewMikrotikConnector(
		host,
		port,
		username,
		password,
		useTls,
		time.Duration(10)*time.Second,
		time.Duration(15)*time.Second,
	)

	for _, item := range collect {
		if metric, ok := item.(map[string]interface{}); ok {
			name := metric["name"].(string)

			//labels := metric["labels"].([]string) //why this doesn't work?
			labelsInterface := metric["labels"].([]interface{})
			labels := make([]string, len(labelsInterface))
			for i, v := range labelsInterface {
				labels[i] = v.(string)
			}

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

			m[name] = NewMikrotikMetric(
				desc,
				valueType,
				connector,
				request,
				context.Background(),
				interval,
			)
		}
	}

	return MikrotikDeviceCollector{
		connector:   connector,
		minInterval: interval,
		metrics:     m,
	}
}

func (c *MikrotikDeviceCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.metrics {
		ch <- prometheus.MustNewConstMetric(
			metric.desc,
			metric.metricType,
			1,                //todo something more clever
			metric.lastValue, // todo match to labels
		)
	}
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
		d.Collect(ch)
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
			d = append(d,
				NewMikrotikDeviceCollector( //isn't it copied?
					targetIP,
					viper.GetInt("targets.port"),
					viper.GetString("targets.username"),
					viper.GetString("targets.password"),
					viper.GetBool("targets.useTls"),
					viper.GetDuration("targets.interval"),
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
