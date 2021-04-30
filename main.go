package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com.sfragata/plex_exporter/collector"
	"github.com.sfragata/plex_exporter/server"
	"github.com/integrii/flaggy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// These variables will be replaced by real values when do gorelease
var (
	version = "none"
	date    string
	commit  string
)

func main() {

	info := fmt.Sprintf(
		"%s\nDate: %s\nCommit: %s\nOS: %s\nArch: %s",
		version,
		date,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)

	flaggy.SetName("plex_exporter")
	flaggy.SetDescription("Prometheus exporter for plex")
	flaggy.SetVersion(info)

	var plexHost = "127.0.0.1"
	flaggy.String(&plexHost, "H", "host", "Plex address")

	var plexPort = 32400
	flaggy.Int(&plexPort, "p", "port", "Plex port")

	var plexToken = os.Getenv("PLEX_TOKEN")
	flaggy.String(&plexToken, "t", "token", "Plex token")

	var metricsPort = "2112"
	flaggy.String(&metricsPort, "l", "listen-address", "Plex exporter metrics port")

	flaggy.Parse()

	if len(plexHost) == 0 {
		flaggy.ShowHelpAndExit("Ples address is mandatory")
	}

	client := &http.Client{
		Timeout: server.HTTPTimeout * time.Second,
	}

	plexServer := server.PlexServer{
		Host:       plexHost,
		Port:       plexPort,
		Token:      plexToken,
		HTTPClient: *client,
	}

	err := prometheus.Register(collector.NewPlexCollector(plexServer))
	if err != nil {
		log.Fatalf("Can't register collectors: %v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("starting plex_export [:%s]", metricsPort)
	err = http.ListenAndServe(":"+metricsPort, nil)
	if err != nil {
		log.Fatalf("Can't start server %s", metricsPort)
	}

}
