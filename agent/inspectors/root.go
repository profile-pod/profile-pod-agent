package inspectors

import (
	"fmt"

	"github.com/profile-pod/profile-pod-agent/agent/details"
	"github.com/profile-pod/profile-pod-agent/agent/inspectors/golang"
	"github.com/profile-pod/profile-pod-agent/agent/profiler"
)

type inspector interface {
	Inspect(process *details.ProcDetails) (profiler.FlameGraphProfiler, bool)
}

var inspectorsList = []inspector{java, python, nodeJs,golang.Go}

// DetectProfiler returns a list of all the detected languages in the process list
func DetectProfiler(p details.ProcDetails) (*profiler.FlameGraphProfiler, error) {
	for _, i := range inspectorsList {
		result, detected := i.Inspect(&p)
		if detected {
			return &result, nil
		}
	}

	return nil, fmt.Errorf("could not find profiler for process %s", p.CmdLine)
}
