package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type fileSinceCollector struct {
	fileSince *prometheus.Desc
}

func newFileSinceCollector() *fileSinceCollector {
	return &fileSinceCollector{
		fileSince: prometheus.NewDesc("file_since_total", "The number of seconds since the file was last modified", []string{"year", "month", "week"}, nil),
	}
}

func (collector *fileSinceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.fileSince
}

func (collector *fileSinceCollector) Collect(ch chan<- prometheus.Metric) {
	FileCounters.mu.RLock()
	defer FileCounters.mu.RUnlock()

	for year, monthToWeekToCount := range FileCounters.data {
		for month, weekToCount := range monthToWeekToCount {
			for week, count := range weekToCount {
				ch <- prometheus.MustNewConstMetric(
					collector.fileSince,
					prometheus.GaugeValue,
					count,
					fmt.Sprintf("%d", year),
					month,
					fmt.Sprintf("%d", week),
				)
			}
		}
	}
}
