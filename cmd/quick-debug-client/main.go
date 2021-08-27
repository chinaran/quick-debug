package main

import (
	"log"
	"os"

	"github.com/chinaran/quick-debug/pkg/client"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "quick-debug",
		Usage: "quick debug program running in the k8s pod",
		Commands: []*cli.Command{
			&uploadCmd,
			&tailLogCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var uploadCmd = cli.Command{
	Name:   "upload",
	Usage:  "uploads a file",
	Action: client.UploadFile,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "addr",
			Usage:    "address of the server to connect to, eg: localhost:60006",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "file",
			Usage:    "the path of file to upload",
			Required: true,
		},
	},
}

var tailLogCmd = cli.Command{
	Name:   "taillog",
	Usage:  "tail log program",
	Action: client.TailLog,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "addr",
			Usage:    "address of the server to connect to, eg: localhost:60006",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "follow",
			Usage: "tail -f, follow the log",
			Value: true,
		},
		&cli.Int64Flag{
			Name:  "n",
			Usage: "tail -n, tail from the last Nth location (byte, not line)",
			Value: 0,
		},
	},
}
