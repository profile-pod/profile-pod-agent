package profiler

import (
	"bytes"
	"os/exec"
	"path"
	"strconv"
	"syscall"

	"github.com/profile-pod/profile-pod-agent/agent/details"
	"github.com/profile-pod/profile-pod-agent/agent/utils/runtime"
	"github.com/profile-pod/profile-pod-agent/agent/utils"
)

const (
	profilerDir = "/tmp/async-profiler"
	fileName    = "/tmp/flamegraph.html"
	profilerSh  = profilerDir + "/profiler.sh"
)

type JvmProfiler struct{}

func (j *JvmProfiler) SetUp(job *details.ProfilingJob) error {
	runtimeFunc, err := runtime.ForRuntime(job.ContainerRuntime)
	targetFs, err := runtimeFunc.GetTargetFileSystemLocation(job.ContainerID)
	if err != nil {
		return err
	}

	err = syscall.Mount(path.Join(targetFs, "tmp"), "/tmp", "", syscall.MS_BIND, "")
	if err != nil {
		return err
	}

	return utils.CopyFolder("/app/async-profiler", "/tmp/async-profiler")
}

func (j *JvmProfiler) Invoke(job *details.ProfilingJob) (string, error) {
	pid := job.ProcDetails.ProcessID

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	event := string(job.Event)
	cmd := exec.Command(profilerSh, "-d", duration, "-f", fileName, "-e", event, pid)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return fileName, nil
}
