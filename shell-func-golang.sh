## 把该文件 shell 函数放入 .zshrc 或 .bashrc 中（或者使用 source /your/script/path 引入）
## 在任意路径直接执行 debug-xxx 类似命令即可
## 可在最下方添加其他程序的调试命令

# if set -e
is_errexit_set() {
  case "$-" in
    *e*) return "true" ;;
    *)   return "false" ;;
  esac
}

# quick-debug
_quick_debug() {
	if [[ $# < 3 ]]; then
		echo "Usage: _quick_debug project-path build-relative-path upload-addr"
		echo "Example: _quick_debug $GOPATH/src/src/github.com/chinaran/go-httpbin ./cmd/go-httpbin 192.168.0.1:32101"
		return
	fi

	if [ "false" == is_errexit_set ]; then
        trap 'set +e' RETURN
        set -e
    fi

	PROJECT_PATH=$1
	BUILD_RELATIVE_PATH=$2
	UPLOAD_ADDR=$3

	TARGET_NAME=${PROJECT_PATH##*/}
	TARGET_PATH=/tmp/$TARGET_NAME
	
	cd $PROJECT_PATH

	# CGO_ENABLED=0 GOOS=linux go build -o /tmp/go-httpbin ./cmd/go-httpbin
	# upx -1 /tmp/go-httpbin
	# quick-debug-client upload --addr 192.168.0.1:32101 --file /tmp/go-httpbin
	# quick-debug-client taillog --addr 192.168.0.1:32101
	echo CGO_ENABLED=0 GOOS=linux go build -o $TARGET_PATH $BUILD_RELATIVE_PATH
	CGO_ENABLED=0 GOOS=linux go build -o $TARGET_PATH $BUILD_RELATIVE_PATH
	upx -1 $TARGET_PATH
	echo quick-debug-client upload --addr $UPLOAD_ADDR --file $TARGET_PATH
	quick-debug-client upload --addr $UPLOAD_ADDR --file $TARGET_PATH
	echo quick-debug-client taillog --addr $UPLOAD_ADDR
	quick-debug-client taillog --addr $UPLOAD_ADDR

	echo
	cd -
}

# TODO(user): add your debug program cmd

# quick-xxx: example shell func, using your real ip and port
debug-xxx() {
	#            project-path                               build-relative-path upload-addr
	_quick_debug $GOPATH/src/github.com/chinaran/go-httpbin ./cmd/go-httpbin 192.168.0.1:32101
}

