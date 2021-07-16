package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func quickDebug() {
	signalCh := make(chan os.Signal)
	quitCh := make(chan struct{})
	// 监听信号
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cmd := new(ExecCmd)
	cmd.Lock()
	go runExec(cmd)
	go func() {
		defer log.Println("exec worker exit")
		for {
			select {
			case <-quitCh:
				if err := cmd.Process.Kill(); err != nil {
					log.Printf("kill %d err: %s", cmd.Process.Pid, err)
				}
				return
			case exec := <-execCh:
				cmd.Lock()
				log.Println("restart with new exec file")
				if err := cmd.Process.Kill(); err != nil {
					log.Fatalf("kill %d err: %s", cmd.Process.Pid, err)
				}
				err := os.Rename(exec.ExecPath, cmdArgs.ExecPath)
				if err != nil {
					log.Fatalf("os.Rename %s to %s err: %s", exec.ExecPath, cmdArgs.ExecPath, err)
				}
				go runExec(cmd)
			}
		}
	}()

	go execFileServer(cmdArgs.ExecPort)

	s := <-signalCh
	log.Println("receive exit signal:", s)
	quitCh <- struct{}{}
}

func runExec(cmd *ExecCmd) {
	cmd.Cmd = exec.Command(cmdArgs.ExecPath, cmdArgs.ExecArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("start run: %s %s", cmdArgs.ExecPath, strings.Join(cmdArgs.ExecArgs, " "))
	err := cmd.Start()
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
