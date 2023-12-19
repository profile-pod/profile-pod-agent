package profiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/utils"
)

const (
	profilerDir = "/tmp/async-profiler"
	fileName    = "/tmp/flamegraph.html"
	profilerSh  = profilerDir + "/profiler.sh"
)

type JvmProfiler struct{}

func (j *JvmProfiler) SetUp(job *details.ProfilingJob) error {
	// runtimeFunc, err := runtime.ForRuntime(job.ContainerRuntime)
	// targetFs, err := runtimeFunc.GetTargetFileSystemLocation(job.ContainerID)
	// if err != nil {
	// 	return err
	// }

	// err = os.RemoveAll("/tmp")
	// if err != nil {
	// 	return err
	// }
	tmpDir := fmt.Sprintf("/proc/%s/root/tmp/", job.ProcDetails.ProcessID)

	err := syscall.Mount(tmpDir, "/tmp", "", syscall.MS_BIND, "")
	if err != nil {
		return err
	}
	//err = os.Symlink(path.Join(targetFs, "tmp"), "/tmp")
	// if err != nil {
	// 	return err
	// }

	return copyFolder("/app/async-profiler", "/tmp/async-profiler")
	//return j.copyProfilerToTempDir()
}

func (j *JvmProfiler) Invoke(job *details.ProfilingJob) error {
	pid  := job.ProcDetails.ProcessID

	duration := strconv.Itoa(int(job.Duration.Seconds()))
	event := string(job.Event)
	cmd := exec.Command(profilerSh, "-d", duration, "-f", fileName, "-e", event, pid)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return utils.PublishFlameGraph(fileName)
}

func (j *JvmProfiler) copyProfilerToTempDir() error {
	cmd := exec.Command("cp", "-r", "/app/async-profiler", "/tmp")
	return cmd.Run()
}

func copyFolder(src, dest string) error {
	// Read the source folder
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Create the destination folder if it doesn't exist
	if err := os.MkdirAll(dest, 0777); err != nil {
		return err
	}

	// Copy each file from source to destination
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			// Recursively copy subdirectories
			if err := copyFolder(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy regular files
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Get source file permissions
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Set destination file permissions
	if err := os.Chmod(dest, srcFileInfo.Mode()); err != nil {
		return err
	}

	return nil
}
