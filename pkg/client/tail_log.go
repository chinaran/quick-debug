package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	pbDebug "github.com/chinaran/quick-debug/proto/debug"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func TailLog(cliCtx *cli.Context) (err error) {
	var (
		addr   = cliCtx.String("addr")
		follow = cliCtx.Bool("follow")
		n      = cliCtx.Int64("n")
	)
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pbDebug.NewQuickDebugClient(conn)
	return tailLog(c, follow, n)
}

func tailLog(c pbDebug.QuickDebugClient, follow bool, n int64) error {
	signalCh := make(chan os.Signal, 2)
	// 监听信号
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	req := pbDebug.TailLogRequest{
		Follow: follow,
		N:      n,
	}
	ctx, cancel := context.WithCancel(context.TODO())
	stream, err := c.TailLog(ctx, &req)
	if err != nil {
		log.Fatalf("failed to call TailLog: %v", err)
	}
	go func() {
		log.Printf("Start show exec running log, follow: %v, n: %v ...\n", follow, n)
		for {
			r, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("failed to Recv server streaming: %v\n", err)
				}
				break
			}
			fmt.Printf("%s\n", r.GetContent())
		}
	}()
	s := <-signalCh
	stream.Context().Done()
	cancel()
	log.Println("receive exit signal:", s)
	return nil
}
