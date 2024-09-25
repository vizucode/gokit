package monitoring

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var once sync.Once

type metrics struct {
	counter *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

var (
	prom        *metrics
	reqsName    = "request_total"
	reqsHelp    = "How many requests processed, partitioned by status code, method, path, and type."
	latencyName = "request_duration_second"
	latencyHelp = "How long it took to process the request, partitioned by status code, method, path, and type."

	DefaultBuckets = []float64{0.3, 1.2, 5.0}
)

func NewPrometheus(serviceName string) {
	once.Do(func() {
		str := []string{"code", "method", "path", "type"}

		reqCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Help:        reqsHelp,
			Name:        reqsName,
			ConstLabels: prometheus.Labels{"service": serviceName},
		}, str)

		if err := prometheus.Register(reqCounter); err != nil {
			return
		}

		reqLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        latencyName,
			Help:        latencyHelp,
			ConstLabels: prometheus.Labels{"service": serviceName},
			Buckets:     DefaultBuckets,
		}, str)

		if err := prometheus.Register(reqLatency); err != nil {
			return
		}

		prom = &metrics{
			counter: reqCounter,
			latency: reqLatency,
		}
	})
}

func PrometheusRecord(statusCode int, method, path, service string, duration time.Duration) {
	if prom == nil {
		return
	}

	code := strconv.Itoa(statusCode)

	var endpoint = path
	if strings.Contains(endpoint, "?") {
		endpoints := strings.Split(endpoint, "?")
		if len(endpoints) > 0 {
			endpoint = endpoints[0]
		}
	}

	prom.counter.WithLabelValues(code, method, endpoint, service).Inc()
	prom.latency.WithLabelValues(code, method, endpoint, service).Observe(float64(duration.Nanoseconds()) / 1000000000)
}
