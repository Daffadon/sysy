package it

import (
	"context"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/daffadon/sysy/test/helper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type GetSyslogITSuite struct {
	suite.Suite
	ctx              context.Context
	network          *testcontainers.DockerNetwork
	gatewayContainer *helper.GatewayContainer
	sysyContainer    *helper.SysyContainer
}

func (g *GetSyslogITSuite) SetupSuite() {
	log.Println("Setting up integration test suite for GetSyslogITSuite")
	g.ctx = context.Background()

	g.network = helper.StartNetwork(g.ctx)
	sysy, err := helper.StartSysyContainer(helper.SysyParameterOption{
		Context:       g.ctx,
		SharedNetwork: g.network.Name,
		ImageName:     "daffaputranarendra/sysy:test",
		ContainerName: "syslog-nginx",
		WaitingSignal: "Exporter listening on :2112/metrics",
		ExposedPorts:  []string{"2112:2112/tcp"},
		Env: map[string]string{
			"CONF_NGINX_TARGET_URL": "http://nginx/nginx_status",
			"CONF_SYSLOG_ADDR":      ":5140",
		},
	})
	if err != nil {
		log.Fatalf("failed starting syslog-nginx container: %s", err)
	}
	g.sysyContainer = sysy

	gatewayContainer, err := helper.StartGatewayContainer(helper.GatewayParameterOption{
		Context:               g.ctx,
		SharedNetwork:         g.network.Name,
		ImageName:             "nginx:1.29.1-alpine-slim",
		ContainerName:         "nginx",
		NginxConfigPath:       "../../config/nginx.conf",
		NginxInsideConfigPath: "/etc/nginx/nginx.conf",
		WaitingSignal:         "Configuration complete; ready for start up",
		MappedPort:            []string{"80:80/tcp"},
	})
	if err != nil {
		log.Fatalf("failed starting gateway container: %s", err)
	}
	g.gatewayContainer = gatewayContainer
	time.Sleep(time.Second)
}

func (g *GetSyslogITSuite) TearDownSuite() {
	if err := g.sysyContainer.Terminate(g.ctx); err != nil {
		log.Fatalf("error terminating sysy container: %s", err)
	}
	if err := g.gatewayContainer.Terminate(g.ctx); err != nil {
		log.Fatalf("error terminating gateway container: %s", err)
	}
	log.Println("Tear Down integration test suite for GetSyslogITSuite")
}

func TestGetSyslogITSuite(t *testing.T) {
	suite.Run(t, &GetSyslogITSuite{})
}

func (g *GetSyslogITSuite) TestGetSyslog_Success() {

	_, err := http.Get("http://localhost/checked")
	g.Require().NoError(err)

	resp, err := http.Get("http://localhost:2112/metrics")
	g.Require().NoError(err)
	defer resp.Body.Close()

	g.Require().Equal(200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	g.Require().NoError(err)

	g.Contains(string(body), "syslog_requests_by_uri")
}
