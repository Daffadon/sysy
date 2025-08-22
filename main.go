package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(requestsByIP, requestsByURI, latencyByURI)
	prometheus.MustRegister(nginxActive, nginxAccepts, nginxHandled, nginxRequests)
	prometheus.MustRegister(nginxReading, nginxWriting, nginxWaiting)
	prometheus.MustRegister(responsesByHTTPCode)
	prometheus.MustRegister(nginxUp)
}

func scrapeStubStatus(t string) {
	resp, err := http.Get(t)
	if err != nil {
		log.Println("Error scraping stub_status:", err)
		nginxUp.Set(0)
		return
	}
	defer resp.Body.Close()
	nginxUp.Set(1)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Active connections:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if v, err := strconv.Atoi(fields[2]); err == nil {
					nginxActive.Set(float64(v))
				}
			}
		} else if strings.Contains(line, "accepts handled requests") {
			continue
		} else if strings.Count(line, " ") >= 2 {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if v, err := strconv.Atoi(fields[0]); err == nil {
					nginxAccepts.Add(float64(v))
				}
				if v, err := strconv.Atoi(fields[1]); err == nil {
					nginxHandled.Add(float64(v))
				}
				if v, err := strconv.Atoi(fields[2]); err == nil {
					nginxRequests.Add(float64(v))
				}
			}
		} else if strings.HasPrefix(line, "Reading:") {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				if v, err := strconv.Atoi(fields[1]); err == nil {
					nginxReading.Set(float64(v))
				}
				if v, err := strconv.Atoi(fields[3]); err == nil {
					nginxWriting.Set(float64(v))
				}
				if v, err := strconv.Atoi(fields[5]); err == nil {
					nginxWaiting.Set(float64(v))
				}
			}
		}
	}
}

func main() {
	_ = godotenv.Load(".env")
	d := os.Getenv("CONF_LOG_ENABLE")
	if d != "true" {
		log.SetOutput(io.Discard)
		fmt.Println("Log Mode Disabled")
	} else {
		fmt.Println("Log Mode Enabled")
	}
	addr := ":5140"
	if v := os.Getenv("CONF_SYSLOG_ADDR"); v != "" {
		addr = v
	}

	go func() {
		pc, err := net.ListenPacket("udp", addr)
		if err != nil {
			panic(err)
		}
		defer pc.Close()

		buf := make([]byte, 65535)
		for {
			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				continue
			}
			line := string(buf[:n])
			if i := strings.Index(line, "nginx:"); i != -1 {
				line = strings.TrimSpace(line[i+7:])
			}
			log.Default().Println(line)
			if entry := parseSyslogLine(line); entry != nil {
				requestsByIP.WithLabelValues(entry.ip).Inc()
				requestsByURI.WithLabelValues(entry.uri).Inc()
				latencyByURI.WithLabelValues(entry.uri).Observe(entry.latency)
				responsesByHTTPCode.WithLabelValues(statusClass(entry.status)).Inc()
			}
		}
	}()

	t := "http://nginx/nginx_status"
	if u := os.Getenv("CONF_NGINX_TARGET_URL"); u != "" {
		t = u
	}

	fmt.Println("Exporter listening on :2112/metrics")
	fmt.Printf("Syslog listening on %s\n", addr)
	fmt.Printf("Nginx Target URL: %s\n", t)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		scrapeStubStatus(t)
		promhttp.Handler().ServeHTTP(w, r)
	})
	http.ListenAndServe(":2112", nil)
}
