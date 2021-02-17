package collector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

const validLibrary3JSON = "{\"MediaContainer\":{\"size\":1}}"
const validLibrary2JSON = "{\"MediaContainer\":{\"size\":2}}"

const validLibrariesJSON = "{\"MediaContainer\":{\"size\":2,\"Directory\":[{\"key\":\"3\",\"type\":\"artist\",\"title\":\"Music\"},{\"key\":\"2\",\"type\":\"movie\",\"title\":\"Children\"}]}}"

func TestLibraryCollectorCollect(test *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), fmt.Sprintf(restLibrary, "3")) {
			fmt.Fprint(w, validLibrary3JSON)
		} else if strings.Contains(r.URL.String(), fmt.Sprintf(restLibrary, "2")) {
			fmt.Fprint(w, validLibrary2JSON)
		} else if strings.Contains(r.URL.String(), restAllLibraries) {
			fmt.Fprint(w, validLibrariesJSON)
		} else {
			test.Errorf("Error: invalid resquest %s", r.URL.String())
		}
	}))
	defer server.Close()

	collector := NewPlexCollector(newPlexServer(server))
	expected := `
	# HELP plex_library_count show number medias
	# TYPE plex_library_count gauge
	plex_library_count{name="Music",type="artist"} 1
	plex_library_count{name="Children",type="movie"} 2
	`
	collector.Collectors = []MetricCollector{LibrariesMetrics{}}

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected)); err != nil {
		test.Errorf("Error:\n%s", err)
	}

}
