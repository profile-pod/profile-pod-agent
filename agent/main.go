// : Copyright Verizon Media
// : Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/inspectors"
	"github.com/VerizonMedia/kubectl-flame/agent/utils"
)

func main() {
	args, err := validateArgs()
	handleError(err, "Error validate Args")

	pDetails, err:= utils.FindRootProcessDetails(args.PodUID,args.ContainerName)
	handleError(err, "Error find process id")

	args.ProcDetails = *pDetails

	p, err := inspectors.DetectProfiler(args.ProcDetails)
	handleError(err, "Error detect profiler")

	err = (*p).SetUp(args)
	handleError(err, "Error SetUp profiler")

	handleSignals()
	flameGraphLocation,err := (*p).Invoke(args);
	handleError(err, "Error Invoke profiler")

	err = utils.PublishFlameGraph(flameGraphLocation)
	handleError(err,"Error Publish Flame Graph")
	cleanUp()
}

func validateArgs() (*details.ProfilingJob, error) {
	if len(os.Args) != 7 && len(os.Args) != 8 {
		return nil, errors.New("expected 7 or 6 arguments")
	}

	duration, err := time.ParseDuration(os.Args[5])
	if err != nil {
		return nil, err
	}

	currentJob := &details.ProfilingJob{}
	currentJob.PodUID = os.Args[1]
	currentJob.ContainerName = os.Args[2]
	currentJob.ContainerID = os.Args[3]
	currentJob.ContainerRuntime= os.Args[4]
	currentJob.Duration = duration
	//currentJob.Language = api.ProgrammingLanguage(os.Args[6])
	currentJob.Event = details.ProfilingEvent(os.Args[6])
	if len(os.Args) == 8 {
		currentJob.TargetProcessName = os.Args[7]
	}

	return currentJob, nil
}

func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		cleanUp()
		//os.Remove("/tmp")
	}()
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf(message + ": %s", err)
		cleanUp()
		os.Exit(1)
	}
}

func cleanUp() {
	os.RemoveAll("/tmp/async-profiler")
}
