package collector

import (
	"log"

	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

//PlexCollector struct that hold metrics collector insterface
type PlexCollector struct {
	Collectors []MetricCollector
	plexServer server.PlexServer
}

//MetricCollector base for the metrics
type MetricCollector interface {
	collect(plexServer server.PlexServer) ([]prometheus.Metric, error)
	describe() *prometheus.Desc
	metricType() prometheus.ValueType
}

//NewPlexCollector constructor
func NewPlexCollector(plexServer server.PlexServer) *PlexCollector {
	return &PlexCollector{

		Collectors: []MetricCollector{SessionMetrics{}, LibrariesMetrics{}, SessionMetricsUserDevice{}},
		plexServer: plexServer,
	}
}

//Describe method from prometheus Collector interface
func (pc *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range pc.Collectors {
		ch <- metric.describe()
	}
}

//Collect method from prometheus Collector interface
func (pc *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range pc.Collectors {

		if metrics, err := collector.collect(pc.plexServer); err == nil {
			for _, metric := range metrics {
				ch <- metric
			}
		} else {
			log.Print(err)
		}
	}
}
