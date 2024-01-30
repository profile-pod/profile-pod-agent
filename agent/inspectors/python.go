package inspectors

import (
	"strings"

	"github.com/profile-pod/profile-pod-agent/agent/details"
	"github.com/profile-pod/profile-pod-agent/agent/profiler"
)

type pythonInspector struct{}

var python = &pythonInspector{}

const pythonProcessName = "python"

func (p *pythonInspector) Inspect(process *details.ProcDetails) (profiler.FlameGraphProfiler, bool) {
	if strings.Contains(process.ExeName, pythonProcessName) || strings.Contains(process.CmdLine, pythonProcessName) {
		return &profiler.PythonProfiler{}, true
	}

	return nil, false
}
