package collector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com.sfragata/plex_exporter/server"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const validSessionJSON = "{\"MediaContainer\":{\"size\":1}}"

func TestSessionCollectorCollect(test *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validSessionJSON)
	}))
	defer server.Close()

	collector := NewPlexCollector(newPlexServer(server))
	expected := `
	# HELP plex_active_sessions_count show number active sessions
	# TYPE plex_active_sessions_count gauge
	plex_active_sessions_count 1
	`
	collector.Collectors = []MetricCollector{SessionMetrics{}}

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		test.Errorf("Error:\n%s", err)
	}

}

func newPlexServer(serverTest *httptest.Server) server.PlexServer {

	hostPort := strings.Split(serverTest.Listener.Addr().String(), ":")

	port, _ := strconv.Atoi(hostPort[1])

	return server.PlexServer{
		Host: hostPort[0],
		Port: port,
	}

}
