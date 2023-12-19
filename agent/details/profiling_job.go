//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package details

import (
	"time"
)

type ProfilingJob struct {
	Duration          time.Duration
	ID                string
	ContainerID       string
	ContainerName     string
	ContainerRuntime  string
	PodUID            string
	TargetProcessName string
	Event             ProfilingEvent
	ProcDetails       ProcDetails
}

type ProcDetails struct {
	ProcessID string
	ExeName   string
	CmdLine   string
}
