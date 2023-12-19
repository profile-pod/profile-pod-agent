package inspectors

import (
	"fmt"
	"os"
	"strings"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/profiler"
)

type javaInspector struct{}

var java = &javaInspector{}

const processName = "java"
const hsperfdataDir = "hsperfdata"

func (j *javaInspector) Inspect(p *details.ProcDetails) (profiler.FlameGraphProfiler, bool) {
	if strings.Contains(p.ExeName, processName) || strings.Contains(p.CmdLine, processName) {
		return &profiler.JvmProfiler{}, true
	}

	if j.searchForHsperfdata(p.ProcessID) {
		return &profiler.JvmProfiler{}, true
	}

	return nil, false
}

func (j *javaInspector) searchForHsperfdata(pid string) bool {
	tmpDir := fmt.Sprintf("/proc/%s/root/tmp/", pid)
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return false
	}

	for _, f := range files {
		if f.IsDir() {
			name := f.Name()
			if strings.Contains(name, hsperfdataDir) {
				return true
			}
		}
	}
	return false
}
