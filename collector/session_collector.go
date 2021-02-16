package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const urlSessions = "http://%s:%d/status/sessions"

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
func (sm SessionMetrics) collectMetrics(ch chan<- prometheus.Metric) error {
	url := fmt.Sprintf(urlSessions, "192.168.2.29", 32400)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	// request.Header.Add("X-Plex-Token", p.Token)
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: status code %d from server", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var responseSession responseSessions

	if err := json.Unmarshal([]byte(body), &responseSession); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(sm.describe(), sm.metricType(), float64(responseSession.MediaContainer.Size))
	return nil
}
