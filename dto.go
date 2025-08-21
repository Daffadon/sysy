package main

type logEntry struct {
	ip      string
	host    string
	method  string
	status  string
	uri     string
	latency float64
}
