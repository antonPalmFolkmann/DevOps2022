package monitoring

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Helge's Example Monitors: A guage representing CPU load %,

var (
	cpuLoad = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "minitwit_cpu_load",
		Help: "The current cpu load percent",
	}, computeCpuLoad)

	responseCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "minitwit_response_counter",
		Help: "The count of total responses sent",
	})

	requestDurationSummary = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "minitwit_request_duration_summary",
		Help:    "The distribution of request durations",
		Buckets: []float64{10.0, 50.0, 100.0, 200.0, 500.0, 1000.0},
	})
)

func init() {
	prometheus.MustRegister(cpuLoad)
	prometheus.MustRegister(responseCounter)
	prometheus.MustRegister(requestDurationSummary)
}

func ResponseCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		responseCounter.Inc()
	})
}

func SetupRoutes(r *mux.Router) {
	r.Use(ResponseCounterMiddleware)

	r.Handle("/metrics", promhttp.Handler())
}

func computeCpuLoad() float64 {
	return 0
}
