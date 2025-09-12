package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/daffadon/sysy/internal/domain/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	ServerReady  chan bool
	Address      string
	ScrapeTarget string
	Logger       *slog.Logger
}

func (m *MetricServer) Run(ctx context.Context) {

	// create slog
	ms := service.NewMetricsService(m.Logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		ms.ScrapeStubStatus(m.ScrapeTarget)
		promhttp.Handler().ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    ":2112",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
	}()

	if m.ServerReady != nil {
		m.ServerReady <- true
	}

	fmt.Printf("Nginx Target URL: %s\n", m.ScrapeTarget)
	fmt.Printf("Exporter listening on %s/metrics\n", m.Address)
	<-ctx.Done()

	server.Shutdown(context.Background())
	fmt.Println("Shutting down metrics server...")
}
