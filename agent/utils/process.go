package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/fntlnz/mountinfo"
)

func getProcessPPID(pid string) (string, error) {
	ppidKey := "PPid"
	statusFile, err := os.Open(path.Join("/proc", pid, "status"))
	if err != nil {
		return "", err
	}

	defer statusFile.Close()
	scanner := bufio.NewScanner(statusFile)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, ppidKey) {
			return strings.Fields(text)[1], nil
		}
	}

	return "", errors.New("unable to get process ppid")
}

func findRootProcess(procsAndParents map[string]string) (string, error) {
	for process, ppid := range procsAndParents {
		if _, ok := procsAndParents[ppid]; !ok {
			// Found process with ppid that is not in the same programming language - this is the root
			return process, nil
		}
	}

	return "", errors.New("could not find root process")
}

func FindRootProcessDetails(podUID string, containerName string) (*details.ProcDetails, error) {
	procsAndParents := make(map[string]details.ProcDetails)
	proc, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}

	for {
		dirs, err := proc.Readdir(15)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, di := range dirs {
			if !di.IsDir() {
				continue
			}

			dname := di.Name()
			if dname[0] < '0' || dname[0] > '9' {
				continue
			}

			//pid, err := strconv.Atoi(dname)
			if err != nil {
				return nil, err
			}

			mi, err := mountinfo.GetMountInfo(path.Join("/proc", dname, "mountinfo"))
			if err != nil {
				continue
			}

			for _, m := range mi {
				root := m.Root
				if strings.Contains(root, fmt.Sprintf("%s/containers/%s", podUID, containerName)) {
					exeName, err := os.Readlink(path.Join("/proc", dname, "exe"))
					if err != nil {
						// Read link may fail if target process runs not as root
						exeName = ""
					}

					cmdLine, err := os.ReadFile(path.Join("/proc", dname, "cmdline"))
					var cmd string
					if err != nil {
						// Ignore errors
						cmd = ""
					} else {
						cmd = string(cmdLine)
					}

					procsAndParents[dname] = details.ProcDetails{
						ProcessID: dname,
						ExeName:   exeName,
						CmdLine:   cmd,
					}
				}
			}
		}
	}

	for process, details := range procsAndParents {
		ppid, err := getProcessPPID(process)
		if err != nil {
			return nil, err
		}
		if _, ok := procsAndParents[ppid]; !ok {
			// Found process with ppid that is not in the same programming language - this is the root
			return &details, nil
		}
	}

	return nil, errors.New("could not find root process")
}
