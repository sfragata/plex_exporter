package collector

import (
	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

type responseSessionsUserDevice struct {
	MediaContainer mediaContainer `json:"MediaContainer"`
}
type mediaContainer struct {
	Size     int        `json:"size"`
	Metadata []metadata `json:"Metadata"`
}
type metadata struct {
	User   user   `json:"User"`
	Player player `json:"Player"`
}

type user struct {
	Title string `json:"title"`
}
type player struct {
	Device string `json:"device"`
}

//SessionMetricsUserDevice SessionMetrics
type SessionMetricsUserDevice struct {
}

func (sm SessionMetricsUserDevice) describe() *prometheus.Desc {
	return prometheus.NewDesc("plex_active_sessions_count_user_device", "show number active sessions by user and device", []string{"user", "device"}, nil)
}
func (sm SessionMetricsUserDevice) metricType() prometheus.ValueType {
	return prometheus.GaugeValue
}
func (sm SessionMetricsUserDevice) collect(plexServer server.PlexServer) ([]prometheus.Metric, error) {
	var responseSession responseSessionsUserDevice
	if err := plexServer.SendRequest(restSessions, &responseSession); err != nil {
		return nil, err
	}

	var metrics []prometheus.Metric
	if responseSession.MediaContainer.Size > 0 {
		for _, metadata := range responseSession.MediaContainer.Metadata {
			metrics = append(metrics, prometheus.MustNewConstMetric(sm.describe(), sm.metricType(), float64(1), metadata.User.Title, metadata.Player.Device))
		}
	}

	return metrics, nil
}
