package collector

import (
	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

const restSessions = "status/sessions"

//ResponseSessions ResponseSessions
type responseSessions struct {
	MediaContainer responseSessionMediaContainer `json:"MediaContainer"`
}

//ResponseSessionMediaContainer ResponseSessionMediaContainer
type responseSessionMediaContainer struct {
	Size int `json:"size"`
}

//SessionMetrics SessionMetrics
type SessionMetrics struct {
}

func (sm SessionMetrics) describe() *prometheus.Desc {
	return prometheus.NewDesc("plex_active_sessions_count", "show number active sessions", nil, nil)
}
func (sm SessionMetrics) metricType() prometheus.ValueType {
	return prometheus.GaugeValue
}
func (sm SessionMetrics) collect(plexServer server.PlexServer) ([]prometheus.Metric, error) {
	var responseSession responseSessions
	if err := plexServer.SendRequest(restSessions, &responseSession); err != nil {
		return nil, err
	}
	metric := prometheus.MustNewConstMetric(sm.describe(), sm.metricType(), float64(responseSession.MediaContainer.Size))

	return []prometheus.Metric{metric}, nil
}
