package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type SysyContainer struct {
	Container testcontainers.Container
}
type SysyParameterOption struct {
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	ExposedPorts                                           []string
	Env                                                    map[string]string
}

func StartSysyContainer(opt SysyParameterOption) (*SysyContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		Env:          opt.Env,
		Networks:     []string{opt.SharedNetwork},
		WaitingFor:   wait.ForLog(opt.WaitingSignal).WithStartupTimeout(30 * time.Second),
		ExposedPorts: opt.ExposedPorts,
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start sysy container: %w", err)
	}

	return &SysyContainer{
		Container: container,
	}, nil
}

func (s *SysyContainer) Terminate(ctx context.Context) error {
	if s.Container != nil {
		return s.Container.Terminate(ctx)
	}
	return nil
}
