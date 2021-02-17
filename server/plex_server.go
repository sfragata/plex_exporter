package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//HTTPTimeout timeout to call endpoints
const HTTPTimeout = 2

const urlTemplate = "http://%s:%d/%s"

//PlexServer plex server information
type PlexServer struct {
	Host       string
	Port       int
	Token      string
	HTTPClient http.Client
}

//SendRequest send requests to plex endpoints
func (ps PlexServer) SendRequest(api string, jsonStruct interface{}) error {

	url := fmt.Sprintf(urlTemplate, ps.Host, ps.Port, api)

	request, _ := http.NewRequest(http.MethodGet, url, nil)

	if len(strings.TrimSpace(ps.Token)) != 0 {
		request.Header.Add("X-Plex-Token", ps.Token)
	}

	request.Header.Add("Accept", "application/json")
	response, err := ps.HTTPClient.Do(request)
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
	defer response.Body.Close()

	if err := json.Unmarshal([]byte(body), &jsonStruct); err != nil {
		return fmt.Errorf("Invalid JSON: %v", err)
	}
	return nil

}
