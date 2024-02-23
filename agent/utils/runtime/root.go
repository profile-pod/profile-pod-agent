package runtime

import (
	"fmt"
)

const (
	dockerRuntime = "docker"
	containerdRuntime = "containerd"
)

type RuntimeFunc interface {
	GetTargetFileSystemLocation(containerId string) (string, error)
}

var (
	docker    = DokcerFunc{}
	containerd    = ContainerdFunc{}
)

func ForRuntime(runtime string) (RuntimeFunc, error) {
	switch runtime {
	case dockerRuntime:
		return &docker, nil
	case containerdRuntime:
		return &containerd, nil
	default:
		return nil, fmt.Errorf("Runtime is not supported: %s",runtime)
	}
}