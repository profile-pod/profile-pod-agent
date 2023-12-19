package inspectors

import (
	"fmt"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/profiler"
)

type inspector interface {
	Inspect(process *details.ProcDetails) (profiler.FlameGraphProfiler, bool)
}

var inspectorsList = []inspector{java, python, nodeJs}

// DetectProfiler returns a list of all the detected languages in the process list
// For go applications the process path is also returned, in all other languages the value is empty
func DetectProfiler(p details.ProcDetails) (*profiler.FlameGraphProfiler, error) {
	for _, i := range inspectorsList {
		result, detected := i.Inspect(&p)
		if detected {
			return &result, nil
		}
	}

	return nil, fmt.Errorf("could not find profiler for process %s", p.CmdLine)
}
