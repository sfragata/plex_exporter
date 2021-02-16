package main

import (
	"net/http"

	"github.com.sfragata/plex_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.Register(collector.NewPlexCollector())
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
