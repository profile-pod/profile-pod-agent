package runtime

import (
	"os"
	"fmt"

	"github.com/profile-pod/profile-pod-agent/agent/common"
)

const (
	containerdMountIdLocation = common.RuntimeBasePath + "/io.containerd.runtime.v2.task/k8s.io/%s/rootfs"
)

type ContainerdFunc struct{}

func (c *ContainerdFunc) GetTargetFileSystemLocation(containerId string) (string, error){
	fileName := fmt.Sprintf(containerdMountIdLocation, containerId)
	_, err := os.Stat(fileName)
	if err == nil {
		// file exists, must be a containerd node
		// the path here is already the rootfs of the container, can return immediately
		return fileName, nil
	}
	return "", err
}