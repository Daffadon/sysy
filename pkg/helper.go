package pkg

import (
	"log"
	"net/url"
	"regexp"
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

var (
	reUUID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	reInt  = regexp.MustCompile(`^\d+$`)
	reHex  = regexp.MustCompile(`(?i)^[0-9a-f]{16,}$`)
)

// NormalizeURI reduces cardinality for Prometheus labels.
func NormalizeURI(raw string) string {
	if raw == "" {
		return "/"
	}

	// If raw is full URL, keep only path.
	if u, err := url.Parse(raw); err == nil && u.Path != "" {
		raw = u.Path
	} else {
		// raw might be "/a/b?x=1"
		if i := strings.IndexByte(raw, '?'); i >= 0 {
			raw = raw[:i]
		}
		if i := strings.IndexByte(raw, '#'); i >= 0 {
			raw = raw[:i]
		}
	}

	// normalize slashes
	path := strings.ReplaceAll(raw, "//", "/")
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	parts := strings.Split(path, "/")
	out := make([]string, 0, len(parts))

	// keep first N segments only
	const maxSegments = 4
	segCount := 0

	for _, p := range parts {
		if p == "" {
			continue
		}
		seg := p

		switch {
		case reUUID.MatchString(seg):
			seg = ":uuid"
		case reInt.MatchString(seg):
			seg = ":id"
		case reHex.MatchString(seg):
			seg = ":hex"
		}

		out = append(out, seg)
		segCount++
		if segCount >= maxSegments {
			out = append(out, ":rest")
			break
		}
	}

	if len(out) == 0 {
		return "/"
	}
	return "/" + strings.Join(out, "/")
}
