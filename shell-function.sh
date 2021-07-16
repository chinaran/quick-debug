## 把该文件 shell 函数放入 .zshrc 或 .bashrc 中（或者使用 source /your/script/path 引入）
## 在任意路径直接执行 debug-xxx 类似命令即可
## 可在最下方添加其他程序的调试命令

# quick-debug
_quick_debug() {
	if [[ $# < 3 ]]; then
		echo "Usage: _quick_debug project-path build-relative-path upload-addr"
		echo "Example: _quick_debug $GOPATH/src/src/github.com/chinaran/go-httpbin ./cmd/go-httpbin 192.168.0.1:32101"
		return
	fi

	PROJECT_PATH=$1
	BUILD_RELATIVE_PATH=$2
	UPLOAD_ADDR=http://$3/upload/exec/file

	TARGET_NAME=${PROJECT_PATH##*/}
	TARGET_PATH=/tmp/$TARGET_NAME
	
	cd $PROJECT_PATH

	# CGO_ENABLED=0 GOOS=linux go build -o /tmp/go-httpbin ./cmd/go-httpbin
	# upx /tmp/go-httpbin
	# curl -X POST http://192.168.0.1:32101/upload/exec/file -F "file=@/tmp/go-httpbin"
	echo CGO_ENABLED=0 GOOS=linux go build -o $TARGET_PATH $BUILD_RELATIVE_PATH
	CGO_ENABLED=0 GOOS=linux go build -o $TARGET_PATH $BUILD_RELATIVE_PATH
	if [[ $? -ne 0 ]];then
		return
	fi
	echo upx $TARGET_PATH
	upx $TARGET_PATH
	echo curl -X POST $UPLOAD_ADDR -F "file=@$TARGET_PATH"
	curl -X POST $UPLOAD_ADDR -F "file=@$TARGET_PATH"

	echo
	cd -
}

# quick-xxx: example shell func
debug-xxx() {
	_quick_debug $GOPATH/src/src/github.com/chinaran/go-httpbin ./cmd/go-httpbin 192.168.0.1:32101
}

