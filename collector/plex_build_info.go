package collector

import (
	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

type responseBuildInfo struct {
	MediaContainer responseBuildInfoMediaContainer `json:"MediaContainer"`
}

//ResponseSessionMediaContainer ResponseSessionMediaContainer
type responseBuildInfoMediaContainer struct {
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	Version         string `json:"version"`
}

//PlexBuildInfoMetrics PlexBuildInfoMetrics
type PlexBuildInfoMetrics struct {
	version string
}

func (pbim PlexBuildInfoMetrics) describe() *prometheus.Desc {
	return prometheus.NewDesc("plex_exporter_build_info", "A metric with a constant '1' value labeled by plex version, plex_exporter version, platform and platformVersion from which plex/plex_exporter was built.", []string{"plex_version", "plex_exporter_version", "platform", "platform_version"}, nil)
}
func (pbim PlexBuildInfoMetrics) metricType() prometheus.ValueType {
	return prometheus.GaugeValue
}
func (pbim PlexBuildInfoMetrics) collect(plexServer server.PlexServer) ([]prometheus.Metric, error) {
	var responseBuildInfo responseBuildInfo
	if err := plexServer.SendRequest("", &responseBuildInfo); err != nil {
		return nil, err
	}
	metric := prometheus.MustNewConstMetric(pbim.describe(), pbim.metricType(), float64(1), responseBuildInfo.MediaContainer.Version, pbim.version, responseBuildInfo.MediaContainer.Platform, responseBuildInfo.MediaContainer.PlatformVersion)

	return []prometheus.Metric{metric}, nil
}
