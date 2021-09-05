# quick-debug

English | [中文](README_ZH.md)

## What Problem To Solve

As the k8s becomes more and more popular, most projects are deployed in k8s, and so is the development environment.

But debugging is a little more troublesome, at least you need to replace the docker image (at this time you also need to commit your code, CICD build image, automatic or manual deployment, which needs much more time).

Timely feedback debugging (similar to front-end local debugging or local development environment) is more friendly to developers and will save a lot of time.

## Application Scenario

* The development environment cannot be built locally, or the local development environment is not real enough
* Quickly debug the program by printing log
* Experiment with unfamiliar lib and check its running results
* Quickly understand the operation of a program
* Develop, debug and self-test at the same time

## Principle

![schematic diagram](https://github.com/chinaran/my-pictures/blob/master/quick-debug/quick-debug-graph.png)

The principle is relatively simple. The image of the development environment needs to contain the quick-debug executable file, and the actual service is started by quick-debug.

After the program is compiled locally (for example, add some debug logs), upload it to the corresponding pod through the exposed node port, restart the service, and achieve the purpose of fast debugging.

## Tool Command

`quick-debug --help`

```shell
Usage: quick-debug [Options] exec-args [real exec args]

Examples:
  quick-debug --exec-port 60006 --exec-path /your/exec/file/path exec-args --arg1 val1 --arg2 val2

Options:
  -disable-exec-log-file
    	disable exec log to file (when disable, you can't using TailLog API)
  -exec-path string
    	exec file path (absolute path)
  -exec-port int
    	exec file server port (Cannot be duplicated with real exec service) (default 60006)
```

`quick-debug-client --help`

```shell
NAME:
   quick-debug - quick debug program running in the k8s pod

USAGE:
   quick-debug-client [global options] command [command options] [arguments...]

COMMANDS:
   upload   uploads a file
   taillog  tail log program
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Effect

### Before Use

Code modification, git submission, CICD generation image (or manual build image), automatic or manual replacement of remote k8s container image, use kubectl or container platform to view logs.

### After Use

You can directly modify the code on the local development machine, compile, run the program on the remote k8s, and view the log on the local machine.

![after use](https://github.com/chinaran/my-pictures/blob/master/quick-debug/result.gif)

## Steps For Usage

Take https://github.com/chinaran/go-httpbin as an example, you can refer to [example](./example/)

**Note:** Your local network needs to be able to access the development k8s cluster, this tool can only be used in the development environment.

### 0. Install quick-debug and quick-debug-client

```shell
git clone https://github.com/chinaran/quick-debug.git $GOPATH/src/github.com/chinaran/quick-debug/
cd $GOPATH/src/github.com/chinaran/quick-debug/
make install
```

### 1. Docker image contains /usr/local/bin/quick-debug

The following methods can be used:

* After local compilation, copy to the image when building
* Take `ghcr.io/chinaran/quick-debug:0.2-alpine3.13` as the base image
* Similar to [docker build example](./example/program-with-quick-debug.Dockerfile), copy `quick-debug` to the image

### 2. Deployment update `command` and `args` 

```yaml
        command:
          - quick-debug
          - --exec-path=/bin/go-httpbin # your exec path in docker
          - exec-args
        args: # orginal args
          - -port
          - "80"
          - -response-delay
          - 10ms
```

The final execution command of go-httpbin pod: `quick-debug --exec-path=/bin/go-httpbin exec-args -port 80 -response-delay 10ms`

The original command before go-httpbin pod modification: `/bin/go-httpbin -port 80 -response-delay 10ms`

### 3. Deploy NodePort Service (default port 60006)

Among them, `nodePort` is best fixed (random ones are fine)

```yaml
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: {your-name}-nodeport
  namespace: {your-namespace}
spec:
  ports:
  - name: quick-deubg
    nodePort: {32101}
    port: 60006
    protocol: TCP
    targetPort: 60006
  selector:
    {your-selector-key}: {your-selector-value}
  type: NodePort
EOF
```

### 4. Compile the executable program locally, use quick-debug-client to upload the program and view its running log

```shell
# Compile the executable file
CGO_ENABLED=0 GOOS=linux go build -o /tmp/go-httpbin ./cmd/go-httpbin
# upx compression
upx -1 /tmp/go-httpbin
# Upload to remote k8s pod
quick-debug-client upload --addr {your-node-ip}:{your-node-port} --file {your-program-path}
# View remote k8s pod running logs locally
quick-debug-client taillog --addr {your-node-ip}:{your-node-port}
```

Golang can use [shell fucntion](./shell-func-golang.sh) to quickly execute corresponding commands (binaries generated in other languages ​​can also refer to this script)
