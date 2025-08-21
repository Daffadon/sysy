package main

import "github.com/prometheus/client_golang/prometheus"

// ----------------- Prometheus Metrics -----------------

var (
	requestsByIP = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_requests_by_ip",
			Help: "Number of requests by client IP",
		},
		[]string{"ip"},
	)

	requestsByURI = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_requests_by_uri",
			Help: "Number of requests by URI",
		},
		[]string{"uri"},
	)

	latencyByURI = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "syslog_request_latency_seconds",
			Help:    "Latency of requests by URI",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"uri"},
	)

	nginxActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_active_connections",
		Help: "Active client connections",
	})

	nginxAccepts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_accepted_connections_total",
		Help: "Total accepted connections",
	})

	nginxHandled = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_handled_connections_total",
		Help: "Total handled connections",
	})

	nginxRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nginx_requests_total",
		Help: "Total requests handled",
	})

	nginxReading = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_reading_connections",
		Help: "Connections where nginx is reading request headers",
	})

	nginxWriting = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_writing_connections",
		Help: "Connections where nginx is writing response back",
	})

	nginxWaiting = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_waiting_connections",
		Help: "Idle connections waiting for requests",
	})

	nginxUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nginx_up",
		Help: "Whether scraping nginx stub_status was successful (1 = UP, 0 = DOWN)",
	})	
	responsesByHTTPCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "syslog_responses_by_class_total",
			Help: "Number of responses grouped by HTTP status class (e.g. 2xx,3xx)",
		},
		[]string{"class"},
	)
)
