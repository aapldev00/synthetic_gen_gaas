package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RecordsGenerated maintains a counter of the total records produced.
	RecordsGenerated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gaas_records_generated_total",
		Help: "The total number of synthetic records generated",
	})

	// ActiveJobs tracks the number of concurrent generation streams.
	ActiveJobs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gaas_active_jobs",
		Help: "Current number of active generation streams",
	})
)

// StartMetricsServer exposes the Prometheus registry over HTTP.
// This allows the scraper to collect telemetry data on a dedicated port.
func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
