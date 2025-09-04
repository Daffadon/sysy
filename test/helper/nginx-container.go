package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GatewayContainer struct {
	Container testcontainers.Container
}
type GatewayParameterOption struct {
	Context                                 context.Context
	SharedNetwork, ImageName, ContainerName string
	NginxConfigPath, NginxInsideConfigPath  string
	WaitingSignal                           string
	MappedPort                              []string
}

func StartGatewayContainer(opt GatewayParameterOption) (*GatewayContainer, error) {
	nginxConfigContent, err := os.ReadFile(opt.NginxConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nginx config file: %w", err)
	}
	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		ExposedPorts: opt.MappedPort,
		WaitingFor:   wait.ForLog(opt.WaitingSignal),
		Networks:     []string{opt.SharedNetwork},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(nginxConfigContent)),
				ContainerFilePath: opt.NginxInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start gateway container: %w", err)

	}

	_, err = container.Host(opt.Context)
	if err != nil {
		container.Terminate(opt.Context)
		return nil, err
	}

	return &GatewayContainer{
		Container: container,
	}, nil
}

func (g *GatewayContainer) Terminate(ctx context.Context) error {
	if g.Container != nil {
		return g.Container.Terminate(ctx)
	}
	return nil
}
