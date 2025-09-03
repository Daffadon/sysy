package variable

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsByIP = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_requests_by_ip",
			Help: "Number of requests by client IP",
		},
		[]string{"ip"},
	)

	RequestsByURI = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_requests_by_uri",
			Help: "Number of requests by URI",
		},
		[]string{"uri"},
	)

	LatencyByURI = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "syslog_request_latency_seconds",
			Help:    "Latency of requests by URI",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"uri"},
	)

	NginxActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_active_connections",
		Help: "Active client connections",
	})

	NginxAccepts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_accepted_connections_total",
		Help: "Total accepted connections",
	})

	NginxHandled = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_handled_connections_total",
		Help: "Total handled connections",
	})

	NginxRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_requests_total",
		Help: "Total requests handled",
	})

	NginxReading = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_reading_connections",
		Help: "Connections where nginx is reading request headers",
	})

	NginxWriting = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_writing_connections",
		Help: "Connections where nginx is writing response back",
	})

	NginxWaiting = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_waiting_connections",
		Help: "Idle connections waiting for requests",
	})

	NginxUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_up",
		Help: "Whether scraping nginx stub_status was successful (1 = UP, 0 = DOWN)",
	})
	ResponsesByHTTPCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_responses_by_class_total",
			Help: "Number of responses grouped by HTTP status class (e.g. 2xx,3xx)",
		},
		[]string{"class"},
	)
)
