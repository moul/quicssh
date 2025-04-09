package main

import (
	"crypto/tls"
	"net"
	"time"

	quic "github.com/quic-go/quic-go"
	cli "github.com/urfave/cli/v2"
)

func client(c *cli.Context) error {
	ctx := withLabel(c.Context, "client")

	config := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quicssh"},
	}

	udpAddr, err := net.ResolveUDPAddr("udp", c.String("addr"))
	if err != nil {
		return er(ctx, err)
	}
	srcAddr, err := net.ResolveUDPAddr("udp", c.String("localaddr"))
	if err != nil {
		return er(ctx, err)
	}

	logf(ctx, "Dialing %q->%q...", srcAddr.String(), udpAddr.String())
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		return er(ctx, err)
	}
	quicConfig := &quic.Config{MaxIdleTimeout: 10 * time.Second, KeepAlivePeriod: 5 * time.Second}
	session, err := quic.Dial(ctx, conn, udpAddr, config, quicConfig)
	if err != nil {
		return er(ctx, err)
	}
	defer func() {
		if err := session.CloseWithError(0, "close"); err != nil {
			logf(ctx, "session close error: %v", err)
		}
	}()

	logf(ctx, "Opening stream sync...")
	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		return er(ctx, err)
	}
	defer stream.Close()

	logf(ctx, "Piping stream with QUIC...")
	c1 := readAndWrite(withLabel(ctx, "stdout"), stream, c.App.Writer) // App.Writer is stdout
	c2 := readAndWrite(withLabel(ctx, "stdin"), c.App.Reader, stream)  // App.Reader is stdin
	select {
	case err = <-c1:
	case err = <-c2:
	}
	if err != nil {
		return err
	}
	return nil
}
