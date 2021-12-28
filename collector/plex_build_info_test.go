package collector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

const validPlexBuildInfoJSON = "{\"MediaContainer\":{\"platform\":\"Linux\",\"platformVersion\":\"5.13.0-22-generic\",\"updatedAt\":1640669424,\"version\":\"1.25.2.5319-c43dc0277\"}}"

func TestPlexBuildInfoCollectorCollect(test *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validPlexBuildInfoJSON)
	}))
	defer server.Close()
	plexExporterVersion := "1.0"
	collector := NewPlexCollector(newPlexServer(server), "")
	expected := fmt.Sprintf(`
	# HELP plex_exporter_build_info A metric with a constant '1' value labeled by plex version, plex_exporter version, platform and platformVersion from which plex/plex_exporter was built.
	# TYPE plex_exporter_build_info gauge
	plex_exporter_build_info{platform="Linux",platform_version="5.13.0-22-generic",plex_exporter_version="%s",plex_version="1.25.2.5319-c43dc0277"} 1
	`, plexExporterVersion)
	collector.Collectors = []MetricCollector{PlexBuildInfoMetrics{version: plexExporterVersion}}

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		test.Errorf("Error:\n%s", err)
	}

}
