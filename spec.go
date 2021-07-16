package main

import (
	"os/exec"
	"sync"
)

var (
	cmdArgs CmdArgs
	execCh  = make(chan *ExecInfo, 10)
)

type CmdArgs struct {
	ExecPort int
	ExecPath string
	ExecArgs []string
	// WeComWebHook string
}

type ExecInfo struct {
	ExecPath string
	// WeComPhone string
}

type ExecCmd struct {
	*exec.Cmd
	sync.Mutex
}
