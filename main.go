package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"

	cli "github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // TODO: application context, good for graceful shutdown
	defer cancel()
	build, _ := debug.ReadBuildInfo()
	app := &cli.App{
		Version: build.Main.Version,
		Usage:   "Client and server parts to proxy SSH (TCP) over UDP using QUIC transport",
		Commands: []*cli.Command{
			{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "bind", Value: "localhost:4242", Usage: "bind address"},
					&cli.StringFlag{Name: "sshdaddr", Value: "localhost:22", Usage: "target address of sshd"},
				},
				Action: server,
			},
			{
				Name: "client",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "addr", Value: "localhost:4242", Usage: "address of server"},
					&cli.StringFlag{Name: "localaddr", Value: ":0", Usage: "source address of UDP packets"},
				},
				Action: client,
			},
		},
	}
	if err := app.RunContext(ctx, os.Args); err != nil {
		logf(ctx, "Error: %v", err)
	}
}

func readAndWrite(ctx context.Context, r io.Reader, w io.Writer) <-chan error {
	c := make(chan error)
	go func() {
		defer close(c)

		buff := make([]byte, 8*1024)

		for {
			select {
			case <-ctx.Done():
				c <- er(ctx, ctx.Err())
				return
			default:
				nr, err := r.Read(buff)
				if err != nil {
					c <- er(ctx, err)
					return
				}
				if nr > 0 {
					_, err := io.Copy(w, bytes.NewReader(buff[:nr]))
					if err != nil {
						c <- er(ctx, err)
						return
					}
				}
			}
		}
	}()
	return c
}

func er(ctx context.Context, e error) error {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Errorf("[%s] %s:%d: %w", label(ctx), path.Base(f), l, e)
}

func logf(ctx context.Context, format string, v ...any) {
	log.Printf("[%s] %s", label(ctx), fmt.Sprintf(format, v...))
}

type lableKeyT int

const lableKey = lableKeyT(0)

func withLabel(ctx context.Context, label string) context.Context {
	if parent, ok := ctx.Value(lableKey).(string); ok {
		label = parent + ">" + label
	}
	return context.WithValue(ctx, lableKey, label)
}

func label(ctx context.Context) string {
	label, _ := ctx.Value(lableKey).(string)
	if label == "" {
		return "main"
	}
	return label
}
