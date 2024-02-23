package profiler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/profile-pod/profile-pod-agent/agent/details"
)

const (
	perfLocation                    = "/app/perf"
	perfRecordOutputFileName        = "/tmp/perf.data"
	flameGraphPlLocation            = "/app/FlameGraph/flamegraph.pl"
	flameGraphStackCollapseLocation = "/app/FlameGraph/stackcollapse-perf.pl"
	perfScriptOutputFileName        = "/tmp/perf.out"
	perfFoldedOutputFileName        = "/tmp/perf.folded"
	perfOutputFilteredFile          = "/tmp/perfFilter.out"
	flameGraphPerfOutputFile        = "/tmp/perf.svg"
)

type PerfProfiler struct{}

func (p *PerfProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (p *PerfProfiler) Invoke(job *details.ProfilingJob) (string, error) {
	err := p.runPerfRecord(job)
	if err != nil {
		return "", fmt.Errorf("perf record failed: %s", err)
	}

	err = p.runPerfScript(job)
	if err != nil {
		return "", fmt.Errorf("perf script failed: %s", err)
	}

	err = p.filterInternalFunction()
	if err != nil {
		return "", fmt.Errorf("filter Internal Function from perf output failed: %s", err)
	}

	err = p.foldPerfOutput(job)
	if err != nil {
		return "", fmt.Errorf("folding perf output failed: %s", err)
	}

	err = p.generateFlameGraph(job)
	if err != nil {
		return "", fmt.Errorf("flamegraph generation failed: %s", err)
	}

	return flameGraphPerfOutputFile, nil
}

func (p *PerfProfiler) runPerfRecord(job *details.ProfilingJob) error {
	duration := strconv.Itoa(int(job.Duration.Seconds()))
	cmd := exec.Command(perfLocation, "record", "-p", job.ProcDetails.ProcessID, "-F99", "-o", perfRecordOutputFileName, "-g", "--", "sleep", duration)

	return cmd.Run()
}

func (p *PerfProfiler) runPerfScript(job *details.ProfilingJob) error {
	f, err := os.Create(perfScriptOutputFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(perfLocation, "script", "-i", perfRecordOutputFileName)
	cmd.Stdout = f

	return cmd.Run()
}

func (p *PerfProfiler) foldPerfOutput(job *details.ProfilingJob) error {
	f, err := os.Create(perfFoldedOutputFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(flameGraphStackCollapseLocation, perfOutputFilteredFile)
	cmd.Stdout = f

	return cmd.Run()
}

func (p *PerfProfiler) filterInternalFunction() error {
	file, err := os.Open(perfScriptOutputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	filteredFile, err := os.Create(perfOutputFilteredFile)
	if err != nil {
		return fmt.Errorf("Error creating filtered file: %s", err)
	}
	defer filteredFile.Close()
	// Create a regular expression for filtering
	filterRegex := regexp.MustCompile(`( __libc_start| LazyCompile | v8::internal::| Builtin:| Stub:| LoadIC:|\[unknown\]| LoadPolymorphicIC:)`)

	// Create a regular expression for substitution
	sedRegex := regexp.MustCompile(` LazyCompile:[*~]?`)

	// Read the file line by line, filter, and substitute
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Apply filtering
		if !filterRegex.MatchString(line) {
			// Apply substitution
			line = sedRegex.ReplaceAllString(line, " ")
			_, err := filteredFile.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("Error writing to filtered file: %s", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: %s", err)
	}

	return nil
}

func (p *PerfProfiler) generateFlameGraph(job *details.ProfilingJob) error {
	inputFile, err := os.Open(perfFoldedOutputFileName)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(flameGraphPerfOutputFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	cmd := exec.Command(flameGraphPlLocation)
	cmd.Stdin = inputFile
	cmd.Stdout = outputFile

	return cmd.Run()
}
