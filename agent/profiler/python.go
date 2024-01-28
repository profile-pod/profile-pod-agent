package profiler

import (
	"bytes"
	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"os/exec"
	"strconv"
)

const (
	pySpyLocation        = "/app/py-spy"
	pythonOutputFileName = "/tmp/python.svg"
)

type PythonProfiler struct{}

func (p *PythonProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (p *PythonProfiler) Invoke(job *details.ProfilingJob) (string,error) {

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(pySpyLocation, "record", "-p", job.ProcDetails.ProcessID, "-o", pythonOutputFileName, "-d", duration, "-s", "-t")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "",err
	}

	return pythonOutputFileName,nil
}
