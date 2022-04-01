package client

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pbDebug "github.com/chinaran/quick-debug/proto/debug"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	ChunkSize = 64 * 1024 // 64KB
)

func UploadFile(cliCtx *cli.Context) (err error) {
	var (
		addr     = cliCtx.String("addr")
		filePath = cliCtx.String("file")
	)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pbDebug.NewQuickDebugClient(conn)
	return uploadFile(c, filePath)
}

func uploadFile(c pbDebug.QuickDebugClient, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("os.Open err: %s\n", err)
		return err
	}
	defer file.Close()

	stream, err := c.UploadFile(context.TODO())
	if err != nil {
		log.Printf("c.UploadFile err: %s\n", err)
		return err
	}
	defer stream.CloseSend()

	startAt := time.Now()
	buf := make([]byte, ChunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("file.Read err: %s\n", err)
			return err
		}
		err = stream.Send(&pbDebug.UploadFileRequest{
			FileChunk: buf[:n],
		})
		if err != nil {
			log.Printf("stream.Send err: %s\n", err)
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("stream.CloseAndRecv err: %s\n", err)
		return err
	}

	log.Printf("Upload %s end, cost: %s\n", filePath, time.Since(startAt))
	log.Printf("Response code: %d, message: %s\n", resp.GetCode(), resp.GetMessage())
	return nil
}
