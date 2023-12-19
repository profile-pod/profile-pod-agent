package profiler

import (
	"github.com/VerizonMedia/kubectl-flame/agent/details"
)

type FlameGraphProfiler interface {
	SetUp(job *details.ProfilingJob) error
	Invoke(job *details.ProfilingJob) error
}