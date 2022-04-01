package server

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// QuickDebug start a quick-debug server
func QuickDebug(cmdArgs *CmdArgs) {
	signalCh := make(chan os.Signal, 2)
	quitCh := make(chan struct{})
	doneCh := make(chan struct{})
	// 监听信号
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cmd := new(ExecCmd)
	cmd.Lock()
	go runExec(cmd, cmdArgs)
	go func() {
		defer log.Println("exec worker exit")
		for {
			select {
			case <-quitCh:
				if err := cmd.Process.Kill(); err != nil {
					log.Printf("kill %d err: %s", cmd.Process.Pid, err)
				}
				doneCh <- struct{}{}
				return
			case exec := <-execCh:
				cmd.Lock()
				log.Println("restart with new exec file")
				if err := cmd.Process.Kill(); err != nil {
					log.Printf("Warn: kill %d err: %s", cmd.Process.Pid, err)
				}
				err := os.Rename(exec.ExecPath, cmdArgs.ExecPath)
				if err != nil {
					log.Fatalf("os.Rename %s to %s err: %s", exec.ExecPath, cmdArgs.ExecPath, err)
				}
				go runExec(cmd, cmdArgs)
			}
		}
	}()

	go startQuickDebugServer(cmdArgs.ExecPort)

	s := <-signalCh
	quitCh <- struct{}{}
	<-doneCh
	log.Println("receive exit signal:", s)
}

func runExec(cmd *ExecCmd, cmdArgs *CmdArgs) {
	var err error
	cmd.Cmd = exec.Command(cmdArgs.ExecPath, cmdArgs.ExecArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if cmdArgs.DisableExecLogFile {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		if cmd.logFile != nil {
			cmd.logFile.Close()
			cmd.logFile = nil
		}
		cmd.logFile, err = os.Create(logFilePath)
		if err != nil {
			log.Fatalf("failed to create %s: %v", logFilePath, err)
		}
		cmd.Stdout = io.MultiWriter(os.Stdout, cmd.logFile)
		cmd.Stderr = io.MultiWriter(os.Stderr, cmd.logFile)
	}
	log.Printf("start run: %s %s", cmdArgs.ExecPath, strings.Join(cmdArgs.ExecArgs, " "))
	err = cmd.Start()
	if err != nil {
		log.Fatalf("failed to call cmd.Start(): %v", err)
	}
	log.Printf("pid: %d", cmd.Process.Pid)
	cmd.Unlock()

	state, err := cmd.Process.Wait()
	if err != nil {
		log.Fatalf("failed to call cmd.Start(): %v", err)
	}
	log.Printf("%s (pid %d) run end, exit info: %s", cmdArgs.ExecPath, state.Pid(), state)
}
