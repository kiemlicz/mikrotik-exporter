package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promslog"
)

var ()

func main() {
	promslogConfig := &promslog.Config{}
	logger := promslog.New(promslogConfig)
	logger.Info("Starting Mikrotik Exporter")

	prometheus.MustRegister(versioncollector.NewCollector("mikrotik_exporter"))

	//srv := &http.Server{}
	//if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
	//	logger.Error("Error starting HTTP server", "err", err)
	//	os.Exit(1)
	//}
}
