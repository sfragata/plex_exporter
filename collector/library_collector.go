package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const urlAllLibraries = "http://%s:%d/library/sections"
const urlLibrary = "http://%s:%d/library/sections/%s/all"

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
func (sm LibrariesMetrics) collectMetrics(ch chan<- prometheus.Metric) error {
	url := fmt.Sprintf(urlAllLibraries, "192.168.2.29", 32400)

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
	// log.Print(string(body))

	var responseLibraries responseLibraries

	if err := json.Unmarshal([]byte(body), &responseLibraries); err != nil {
		return err
	}
	for _, directory := range responseLibraries.MediaContainer.ResponseDirectories {
		url := fmt.Sprintf(urlLibrary, "192.168.2.29", 32400, directory.Key)

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
		// log.Print(string(body))
		if err != nil {
			return err
		}
		var responseLibrary responseLibrary

		if err := json.Unmarshal([]byte(body), &responseLibrary); err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(sm.describe(), sm.metricType(), float64(responseLibrary.MediaContainer.Size), directory.Title, directory.Type)
	}

	return nil
}
