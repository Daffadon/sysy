package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	v "github.com/daffadon/sysy/internal/variable"
	"github.com/daffadon/sysy/pkg"
)

type SysServer struct {
	ServerReady chan bool
	Address     string
}

func (s *SysServer) Run(ctx context.Context) {
	addr := s.Address
	go func() {
		pc, err := net.ListenPacket("udp", addr)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Syslog listening on %s\n", addr)

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
			if entry := pkg.ParseSyslogLine(line); entry != nil {
				v.RequestsByIP.WithLabelValues(entry.IP).Inc()
				v.RequestsByURI.WithLabelValues(entry.URI).Inc()
				v.LatencyByURI.WithLabelValues(entry.URI).Observe(entry.Latency)
				v.ResponsesByHTTPCode.WithLabelValues(pkg.StatusClass(entry.Status)).Inc()
			}
		}
	}()
	if s.ServerReady != nil {
		s.ServerReady <- true
	}
	<-ctx.Done()
	log.Default().Println("shutting down syslog server")
}
