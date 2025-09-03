package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/daffadon/sysy/cmd/server"
	v "github.com/daffadon/sysy/internal/variable"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(v.RequestsByIP, v.RequestsByURI, v.LatencyByURI)
	prometheus.MustRegister(v.NginxActive, v.NginxAccepts, v.NginxHandled, v.NginxRequests)
	prometheus.MustRegister(v.NginxReading, v.NginxWriting, v.NginxWaiting)
	prometheus.MustRegister(v.ResponsesByHTTPCode)
	prometheus.MustRegister(v.NginxUp)
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
	t := "http://nginx/nginx_status"
	if u := os.Getenv("CONF_NGINX_TARGET_URL"); u != "" {
		t = u
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsServerReady := make(chan bool)
	metricsServerDone := make(chan struct{})
	metricsServer := &server.MetricServer{
		ServerReady:  metricsServerReady,
		Address:      ":2112",
		ScrapeTarget: t,
	}
	go func() {
		metricsServer.Run(ctx)
		close(metricsServerDone)
	}()

	<-metricsServerReady

	addr := ":5140"
	if v := os.Getenv("CONF_SYSLOG_ADDR"); v != "" {
		addr = v
	}

	sylogServerReady := make(chan bool)
	sylogServerDone := make(chan struct{})
	sylogServer := &server.SysServer{
		ServerReady: sylogServerReady,
		Address:     addr,
	}
	go func() {
		sylogServer.Run(ctx)
		close(sylogServerDone)
	}()

	<-sylogServerReady

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)

	<-sig
	cancel()

	<-metricsServerDone
	<-sylogServerDone
}
