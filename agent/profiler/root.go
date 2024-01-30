package profiler

import "github.com/profile-pod/profile-pod-agent/agent/details"

type FlameGraphProfiler interface {
	SetUp(job *details.ProfilingJob) error
	Invoke(job *details.ProfilingJob) (string, error)
}
