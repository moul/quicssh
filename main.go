package main

import (
	"io"
	"os"
	"sync"

	"golang.org/x/net/context"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "bind", Value: "localhost:4242"},
				},
				Action: server,
			},
			&cli.Command{
				Name: "client",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "addr", Value: "localhost:4242"},
				},
				Action: client,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func readAndWrite(ctx context.Context, r io.Reader, w io.Writer, wg *sync.WaitGroup) <-chan error {
	c := make(chan error)
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		buff := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				nr, err := r.Read(buff)
				if err != nil {
					return
				}
				if nr > 0 {
					_, err := w.Write(buff[:nr])
					if err != nil {
						return
					}
				}
			}
		}
	}()
	return c
}
