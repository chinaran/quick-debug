# quick-debug

快速调试运行在 k8s pod 中的程序 

过程：本地编译好可执行文件，上传到对应的 pod，quick-debug 重启目标程序

## 使用步骤

以 https://github.com/chinaran/go-httpbin 为例

**注:** 本机需要能访问开发 k8s 集群

### 1. docker image 中 包含 /usr/local/bin/quick-debug

### 2. deployment 更新 `command` 和 `args` 

```yaml
        command:
          - quick-debug
          - --exec-path=/bin/go-httpbin
          - exec-args
        args:
          - -port
          - 80
          - -response-delay
          - 10ms
```

go-httpbin pod 最终执行命令: `quick-debug --exec-path=/bin/go-httpbin exec-args -port 80 -response-delay 10ms`

### 3. 部署 NodePort Service (默认端口 60006)

其中 `nodePort` 最好固定好（使用随机的也行）

```yaml
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
```

### 4. 本地编译好可执行程序，使用 quick-debug-client 上传程序和查看其运行日志

```shell
CGO_ENABLED=0 GOOS=linux go build -o /tmp/go-httpbin ./cmd/go-httpbin
upx /tmp/go-httpbin
quick-debug-client upload --addr {your-node-ip}:{your-node-port} --file {your-program-path}
quick-debug-client taillog --addr {your-node-ip}:{your-node-port}
```
