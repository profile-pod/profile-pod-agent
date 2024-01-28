package profiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
)

const (
	rbspyLocation       = "/app/rbspy"
	rbspyOutputFileName = "/tmp/rbspy"
)

type RubyProfiler struct{}

func (r *RubyProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (r *RubyProfiler) Invoke(job *details.ProfilingJob) (string,error) {

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(rbspyLocation, "record", "--pid", job.ProcDetails.ProcessID, "--file", rbspyOutputFileName, "--duration", duration, "--format", "flamegraph")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "",fmt.Errorf("could not launch profiler: %w", err)
	}

	return rbspyOutputFileName,nil
}
