//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/profiler"
	"github.com/VerizonMedia/kubectl-flame/api"
)

func main() {
	args, err := validateArgs()
	handleError(err)

	// err = api.PublishEvent(api.Progress, &api.ProgressData{Time: time.Now(), Stage: api.Started})
	// handleError(err)

	p, err := profiler.ForLanguage(args.Language)
	handleError(err)

	err = p.SetUp(args)
	handleError(err)

	handleSignals()
	err = p.Invoke(args)
	handleError(err)
	cleanUp()
	// err = api.PublishEvent(api.Progress, &api.ProgressData{Time: time.Now(), Stage: api.Ended})
	// handleError(err)
}

func validateArgs() (*details.ProfilingJob, error) {
	if len(os.Args) != 8 && len(os.Args) != 9 {
		return nil, errors.New("expected 7 or 8 arguments")
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
	currentJob.Language = api.ProgrammingLanguage(os.Args[6])
	currentJob.Event = api.ProfilingEvent(os.Args[7])
	if len(os.Args) == 9 {
		currentJob.TargetProcessName = os.Args[8]
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

func handleError(err error) {
	if err != nil {
		api.PublishError(err)
		cleanUp()
		os.Exit(1)
	}
}

func cleanUp() {
	os.RemoveAll("/tmp/async-profiler")
}
