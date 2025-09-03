package pkg

import (
	"log"
	"strconv"
	"strings"

	"github.com/daffadon/sysy/internal/domain/dto"
)

func ParseSyslogLine(line string) *dto.LogEntry {
	parts := strings.Split(line, "\t")
	if len(parts) < 8 {
		log.Printf("Error Log Parse: %s\n", line)
		return nil
	}

	lat, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		lat = 0
	}

	return &dto.LogEntry{
		IP:      parts[0],
		Method:  parts[1],
		Status:  parts[2],
		URI:     parts[7],
		Latency: lat,
	}
}

func StatusClass(status string) string {
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
