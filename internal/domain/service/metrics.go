package service

import (
	"bufio"
	"log"
	"net/http"
	"strconv"
	"strings"

	vars "github.com/daffadon/sysy/internal/variable"
)

type (
	MetricsService interface {
		ScrapeStatus(t string)
	}
)

func ScrapeStubStatus(t string) {
	resp, err := http.Get(t)
	if err != nil {
		log.Println("Error scraping stub_status:", err)
		vars.NginxUp.Set(0)
		return
	}
	defer resp.Body.Close()
	vars.NginxUp.Set(1)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Active connections:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if v, err := strconv.Atoi(fields[2]); err == nil {
					vars.NginxActive.Set(float64(v))
				}
			}
		} else if strings.Contains(line, "accepts handled requests") {
			continue
		} else if strings.Count(line, " ") >= 2 {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if v, err := strconv.Atoi(fields[0]); err == nil {
					vars.NginxAccepts.Add(float64(v))
				}
				if v, err := strconv.Atoi(fields[1]); err == nil {
					vars.NginxHandled.Add(float64(v))
				}
				if v, err := strconv.Atoi(fields[2]); err == nil {
					vars.NginxRequests.Add(float64(v))
				}
			}
		} else if strings.HasPrefix(line, "Reading:") {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				if v, err := strconv.Atoi(fields[1]); err == nil {
					vars.NginxReading.Set(float64(v))
				}
				if v, err := strconv.Atoi(fields[3]); err == nil {
					vars.NginxWriting.Set(float64(v))
				}
				if v, err := strconv.Atoi(fields[5]); err == nil {
					vars.NginxWaiting.Set(float64(v))
				}
			}
		}
	}
}
