package collector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

const validSessionUserDeviceJSON = "{\"MediaContainer\":{\"size\":2,\"Metadata\":[{\"User\":{\"title\":\"user1\"},\"Player\":{\"device\":\"iPad\"}},{\"User\":{\"title\":\"user2\"},\"Player\":{\"device\":\"iPhone\"}}]}}"

func TestSessionUserDeviceCollectorCollect(test *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, validSessionUserDeviceJSON)
	}))
	defer server.Close()

	collector := NewPlexCollector(newPlexServer(server))
	expected := `
	# HELP plex_active_sessions_count_user_device show number active sessions by user and device
	# TYPE plex_active_sessions_count_user_device gauge
	plex_active_sessions_count_user_device{device="iPad",user="user1"} 1
	plex_active_sessions_count_user_device{device="iPhone",user="user2"} 1
	`
	collector.Collectors = []MetricCollector{SessionMetricsUserDevice{}}

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		test.Errorf("Error:\n%s", err)
	}

}
