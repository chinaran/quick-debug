package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pbDebug "github.com/chinaran/quick-debug/proto/debug"
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

type QuickDebugServer struct {
	pbDebug.UnimplementedQuickDebugServer
}

func NewQuickDebugServer() *QuickDebugServer {
	return &QuickDebugServer{}
}

func startQuickDebugServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Quick-Debug grpc server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pbDebug.RegisterQuickDebugServer(s, NewQuickDebugServer())
	reflection.Register(s)
	s.Serve(lis)
}

func (s *QuickDebugServer) UploadFile(stream pbDebug.QuickDebug_UploadFileServer) error {
	log.Println("Start receive file ...")
	// Save the file to specific dst path.
	dstPath := fmt.Sprintf("/tmp/%s", uuid.NewString())
	f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("os.OpenFile err: %s", err)
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// receive left file chunk
				if len(chunk.GetFileChunk()) > 0 {
					if _, err := f.Write(chunk.FileChunk); err != nil {
						log.Printf("f.Write err: %s", err)
						f.Close()
						return err
					}
				}
				f.Close()
				log.Println("Finish receive file")

				execCh <- &ExecInfo{ExecPath: dstPath}

				err = stream.SendAndClose(&pbDebug.UploadFileResponse{
					Code:    int32(codes.OK),
					Message: "OK",
				})
				if err != nil {
					log.Printf("stream.SendAndClose err: %s", err)
					return err
				}
				return nil
			}
			log.Printf("stream.Recv err: %s", err)
			f.Close()
			return err
		}
		// receive file chunk
		if len(chunk.GetFileChunk()) > 0 {
			if _, err := f.Write(chunk.FileChunk); err != nil {
				log.Printf("f.Write err: %s", err)
				f.Close()
				return err
			}
		}
	}
}

func (*QuickDebugServer) TailLog(req *pbDebug.TailLogRequest, stream pbDebug.QuickDebug_TailLogServer) error {
	config := tail.Config{
		Follow: req.GetFollow(),
	}
	if req.GetN() != 0 {
		config.Location = &tail.SeekInfo{Offset: -req.N, Whence: io.SeekEnd}
	}
	t, err := tail.TailFile(logFilePath, config)
	if err != nil {
		log.Printf("tail.TailFile err: %s", err)
		return err
	}
	defer t.Stop()
	for {
		select {
		case <-stream.Context().Done():
			log.Println("TailLog end by stream.Context().Done()")
			return nil
		case line := <-t.Lines:
			err := stream.Send(&pbDebug.TailLogResponse{Content: line.Text})
			if err != nil {
				log.Printf("stream.Send err: %s", err)
				return err
			}
		}
	}
}
