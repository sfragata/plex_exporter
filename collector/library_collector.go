package collector

import (
	"fmt"

	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

const restAllLibraries = "library/sections"
const restLibrary = "library/sections/%s/all"

//ResponseLibrary ResponseLibrary
type responseLibrary struct {
	MediaContainer responseLibrarySize `json:"MediaContainer"`
}

//ResponseLibrarySize ResponseLibrarySize
type responseLibrarySize struct {
	Size int `json:"size"`
}

type responseLibraries struct {
	MediaContainer responseLibraryDirectories `json:"MediaContainer"`
}

type responseLibraryDirectories struct {
	ResponseDirectories []responseDirectory `json:"Directory"`
}

type responseDirectory struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

//LibrariesMetrics LibrariesMetrics
type LibrariesMetrics struct {
}

func (sm LibrariesMetrics) describe() *prometheus.Desc {
	return prometheus.NewDesc("plex_library_count", "show number medias", []string{"name", "type"}, nil)
}
func (sm LibrariesMetrics) metricType() prometheus.ValueType {
	return prometheus.GaugeValue
}
func (sm LibrariesMetrics) collect(plexServer server.PlexServer) ([]prometheus.Metric, error) {

	var responseLibraries responseLibraries

	if err := plexServer.SendRequest(restAllLibraries, &responseLibraries); err != nil {
		return nil, err
	}

	var metrics []prometheus.Metric
	for _, directory := range responseLibraries.MediaContainer.ResponseDirectories {

		restEachLibrary := fmt.Sprintf(restLibrary, directory.Key)

		var responseLibrary responseLibrary

		if err := plexServer.SendRequest(restEachLibrary, &responseLibrary); err != nil {
			return nil, err
		}

		metrics = append(metrics, prometheus.MustNewConstMetric(sm.describe(), sm.metricType(), float64(responseLibrary.MediaContainer.Size), directory.Title, directory.Type))

	}
	return metrics, nil
}
