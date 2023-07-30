package main

import (
	"crypto/tls"
	"log"
	"os"
	"sync"

	quic "github.com/quic-go/quic-go"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/net/context"
)

func client(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	config := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quicssh"},
	}

	log.Printf("Dialing %q...", c.String("addr"))
	session, err := quic.DialAddr(ctx, c.String("addr"), config, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := session.CloseWithError(0, "close"); err != nil {
			log.Printf("session close error: %v", err)
		}
	}()

	log.Printf("Opening stream sync...")
	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		return err
	}

	log.Printf("Piping stream with QUIC...")
	var wg sync.WaitGroup
	wg.Add(3)
	c1 := readAndWrite(ctx, stream, os.Stdout, &wg)
	c2 := readAndWrite(ctx, os.Stdin, stream, &wg)
	select {
	case err = <-c1:
		if err != nil {
			return err
		}
	case err = <-c2:
		if err != nil {
			return err
		}
	}
	cancel()
	wg.Wait()
	return nil
}
