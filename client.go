package main

import (
	"crypto/tls"
	"log"
	"os"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
	"golang.org/x/net/context"
	cli "gopkg.in/urfave/cli.v2"
)

func client(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	config := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quicssh"},
	}

	log.Printf("Dialing %q...", c.String("addr"))
	session, err := quic.DialAddr(c.String("addr"), config, nil)
	if err != nil {
		return err
	}
	defer session.Close()

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
