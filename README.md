# quick-debug

## 要解决什么问题

随着 k8s 的越来越流行，大多数项目都部署在 k8s 中，开发环境也是。

但调试稍微麻烦，至少需要替换 docker 镜像（此时还需要提 git 提交，CICD 构建镜像，自动或手动部署，时间较长）。

可见即所得的调试（类似前端本地调试或本地有开发环境）对开发者更友好，也会节约很多时间。

## 应用场景

* 本地无法搭建开发环境，或本地开发环境不够真实
* 通过打印日志方式快速调试程序，debug
* 试验不熟悉的 lib，查看其运行结果
* 快速了解一个程序的运行情况
* 边开发，边调试，边自测

## 原理

![schematic diagram](https://github.com/chinaran/my-pictures/blob/master/quick-debug/quick-debug-graph.png)

原理比较简单，开发环境的镜像需要包含 quick-debug 可执行文件，实际的服务由 quick-debug 启动。

当本地编译好程序后（例如加一些 debug 日志），通过暴露的 node port 上传到对应 pod，重启服务，达到快速调试的目的。

## 工具命令

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

## 使用效果

### 使用前

代码修改，git 提交，CICD 生成镜像（或手动 build 镜像），自动或手动替换远端 k8s 容器镜像，使用 kubectl 或容器平台查看日志

### 使用后

可直接在本地开发机修改代码，编译，在远端 k8s 运行程序，在本机查看日志

![after use](https://github.com/chinaran/my-pictures/blob/master/quick-debug/result.gif)

## 使用步骤

以 https://github.com/chinaran/go-httpbin 为例，可参考 [example](./example/)

**注:** 本机网络需要能访问开发 k8s 集群，该工具只可用于开发环境

### 0. 安装 quick-debug 和 quick-debug-client

```shell
git clone https://github.com/chinaran/quick-debug.git $GOPATH/src/github.com/chinaran/quick-debug/
cd $GOPATH/src/github.com/chinaran/quick-debug/
make install
```

### 1. docker image 中 包含 /usr/local/bin/quick-debug

可使用如下方式：

* 本地编译后，构建时复制到镜像中
* 以 `ghcr.io/chinaran/quick-debug:0.2-alpine3.13` 作为基础镜像
* 类似 [docker build example](./example/program-with-quick-debug.Dockerfile)，复制 `quick-debug` 到镜像中

### 2. deployment 更新 `command` 和 `args` 

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

go-httpbin pod 最终执行命令: `quick-debug --exec-path=/bin/go-httpbin exec-args -port 80 -response-delay 10ms`

go-httpbin pod 修改前原命令: `/bin/go-httpbin -port 80 -response-delay 10ms`

### 3. 部署 NodePort Service (默认端口 60006)

其中 `nodePort` 最好固定好（使用随机的也行）

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

### 4. 本地编译好可执行程序，使用 quick-debug-client 上传程序和查看其运行日志

```shell
# 编译可执行文件
CGO_ENABLED=0 GOOS=linux go build -o /tmp/go-httpbin ./cmd/go-httpbin
# upx 压缩
upx -1 /tmp/go-httpbin
# 上传到远端 k8s pod
quick-debug-client upload --addr {your-node-ip}:{your-node-port} --file {your-program-path}
# 在本地查看远端 k8s pod 运行日志
quick-debug-client taillog --addr {your-node-ip}:{your-node-port}
```

Golang 可使用 [shell fucntion](./shell-func-golang.sh) 快速执行对应命令 (其他语言生成的二进制也可参考该脚本)
