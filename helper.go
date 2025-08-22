package main

import (
	"log"
	"strconv"
	"strings"
)

func parseSyslogLine(line string) *logEntry {
	parts := strings.Split(line, "\t")
	if len(parts) < 8 {
		log.Printf("Error Log Parse: %s\n", line)
		return nil
	}

	lat, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		lat = 0
	}

	return &logEntry{
		ip:      parts[0],
		method:  parts[1],
		status:  parts[2],
		uri:     parts[7],
		latency: lat,
	}
}

// statusClass converts a status string like "200" to its class label "2xx".
func statusClass(status string) string {
	if len(status) >= 3 {
		switch status[0] {
		case '2':
			return "2xx"
		case '3':
			return "3xx"
		case '4':
			return "4xx"
		case '5':
			return "5xx"
		}
	}
	return "other"
}
