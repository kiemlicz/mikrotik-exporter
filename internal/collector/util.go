package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

func parseValueType(t string) (prometheus.ValueType, error) {
	switch t {
	case "gauge":
		return prometheus.GaugeValue, nil
	case "counter":
		return prometheus.CounterValue, nil
	default:
		return 0, fmt.Errorf("unknown metric type: %s", t)
	}
}
