package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var configuration = &Configuration{}

func main() {
	configuration.Parse()

	ticker := time.NewTicker(configuration.WalkingInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			doWalk(configuration)
		}
	}()

	prometheus.MustRegister(newFileSinceCollector())

	http.Handle("/metrics", promhttp.Handler())

	log.Info("Listening on ", configuration.ListenAddress)
	log.Fatalln(http.ListenAndServe(configuration.ListenAddress, nil))
}
