package collector

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const httpTimeout = 2

//PlexCollector struct that hold metrics collector insterface
type PlexCollector struct {
	metrics []MetricCollector
}

//MetricCollector base for the metrics
type MetricCollector interface {
	collectMetrics(ch chan<- prometheus.Metric) error
	describe() *prometheus.Desc
	metricType() prometheus.ValueType
}

var (
	client = &http.Client{
		Timeout: httpTimeout * time.Second,
	}
)

//NewPlexCollector constructor
func NewPlexCollector() *PlexCollector {
	return &PlexCollector{

		metrics: []MetricCollector{SessionMetrics{}, LibrariesMetrics{}},
	}
}

//Describe method from prometheus Collector interface
func (pc *PlexCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range pc.metrics {
		ch <- metric.describe()
	}
}

//Collect method from prometheus Collector interface
func (pc *PlexCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range pc.metrics {

		if err := metric.collectMetrics(ch); err != nil {
			log.Print(err)
		}
	}
}
