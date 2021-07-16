package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// curl -X POST http://localhost:60006/upload/exec/file -F "file=@/your/file/path"
func execFileServer(port int) {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.POST("/upload/exec/file", uploadExecFile)

	router.Run(fmt.Sprintf(":%d", port))
}

func uploadExecFile(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		ginResp(c, http.StatusBadRequest, fmt.Sprintf("get file err: %s", err))
		return
	}

	// Upload the file to specific dst.
	dst := fmt.Sprintf("/tmp/%s", uuid.NewString())
	c.SaveUploadedFile(file, dst)
	err = os.Chmod(dst, 0755)
	if err != nil {
		ginResp(c, http.StatusInternalServerError, fmt.Sprintf("os.Chmod %s err: %s", dst, err))
		return
	}

	execCh <- &ExecInfo{ExecPath: dst}

	ginResp(c, http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func ginResp(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
