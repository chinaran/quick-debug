package server

import (
	"os"
	"os/exec"
	"sync"
)

const (
	logFilePath = "/tmp/quick-debug-exec.log"
)

var (
	execCh = make(chan *ExecInfo, 10)
)

// CmdArgs ...
type CmdArgs struct {
	ExecPort           int
	ExecPath           string
	DisableExecLogFile bool
	ExecArgs           []string
}

// ExecInfo ...
type ExecInfo struct {
	ExecPath string
}

// ExecCmd ...
type ExecCmd struct {
	*exec.Cmd
	sync.Mutex
	logFile *os.File
}
