package runtime

import (
	"os"
	"fmt"

	"github.com/profile-pod/profile-pod-agent/agent/common"
)

const (
	overlayFsMountIdLocation = common.RuntimeBasePath + "/image/overlay2/layerdb/mounts/%s/mount-id"
	targetFileSystemLocation  = common.RuntimeBasePath + "/overlay2/%s/merged"
)

type DokcerFunc struct{}

func (d *DokcerFunc) GetTargetFileSystemLocation(containerId string) (string, error){
	fileName := fmt.Sprintf(overlayFsMountIdLocation, containerId)
	mountId, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(targetFileSystemLocation, string(mountId)), nil
}