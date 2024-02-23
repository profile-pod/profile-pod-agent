package inspectors

import (
	"strings"

	"github.com/profile-pod/profile-pod-agent/agent/details"
	"github.com/profile-pod/profile-pod-agent/agent/profiler"
)

type nodejsInspector struct{}

var nodeJs = &nodejsInspector{}

const nodeProcessName = "node"

func (n *nodejsInspector) Inspect(process *details.ProcDetails) (profiler.FlameGraphProfiler, bool) {
	if strings.Contains(process.ExeName, nodeProcessName) || strings.Contains(process.CmdLine, nodeProcessName) {
		return &profiler.PerfProfiler{}, true
	}

	return nil, false
}
