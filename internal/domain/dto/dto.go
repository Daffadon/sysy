package dto

type LogEntry struct {
	IP      string
	Method  string
	Status  string
	URI     string
	Latency float64
}
