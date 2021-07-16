package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.IntVar(&cmdArgs.ExecPort, "exec-port", 60006, "exec file server port (Cannot be duplicated with real exec service)")
	flag.StringVar(&cmdArgs.ExecPath, "exec-path", "", "exec file path (absolute path)")
	// flag.StringVar(&cmdArgs.WeComWebHook, "wecom-webhook", "", "Send WeCom robot message")

	flag.Usage = usage
}

func usage() {
	fmt.Printf(`Usage: quick-debug [Options] exec-args [real exec args]

Examples:
  quick-debug --exec-port 60006 --exec-path /your/exec/file/path exec-args --arg1 val1 --arg2 val2

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if cmdArgs.ExecPath == "" {
		fmt.Println(`Required flag "exec-path" not set`)
		return
	}
	if len(flag.Args()) > 0 {
		cmdArgs.ExecArgs = flag.Args()[1:]
	}
	quickDebug()
}
