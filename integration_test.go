package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"net"
	"strings"
	"testing"

	cli "github.com/urfave/cli/v2"
)

const (
	sshdAddr = "localhost:8822"
	bindAddr = "localhost:4242"
)

func TestInegratinGoldenFlow(t *testing.T) {
	// run tcp-echo server (sort of sshd)
	closer := tcpEchoServer()
	defer closer() // close fake sshd

	// run quicssh server
	serverContext := flags("sshdaddr", sshdAddr, "bind", bindAddr)
	go func() {
		err := server(serverContext) // TODO: make it able to shutdown gracefully, check error
		if strings.Contains(err.Error(), "address already in use") {
			panic(err.Error()) // cheapest way to fail test from goroutine
		}
	}()

	// run quicssh client
	clientContext := flags("addr", bindAddr, "localaddr", ":0")
	wr, rd := tweaksStdIO(clientContext)
	go func() {
		err := client(clientContext) // TODO: shutdownable, check error
		if !errors.Is(err, io.EOF) {
			panic(err.Error())
		}
	}()

	// writing data to client
	_, err := wr.Write([]byte("SSH SESSION\n")) // writing one line
	noerr(err)

	// reading echo-data from client<-proxyserver<-echotcpserver<-proxyserver
	lineReader := bufio.NewReader(rd)
	line, err := lineReader.ReadString('\n') // waiting for just one echo-line back
	noerr(err)

	// doing checks
	const expected = "echo: SSH SESSION\n"
	if line != expected {
		t.Fatalf("expected %q, got %q", expected, line)
	}
}

func tweaksStdIO(c *cli.Context) (io.Writer, io.Reader) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c.App.Reader = r1
	c.App.Writer = w2
	return w1, r2 // returning the opposite end of pipes
}

func flags(pairs ...string) *cli.Context {
	fs := flag.NewFlagSet("", 0)
	for i := 0; i < len(pairs); i += 2 {
		fs.String(pairs[i], pairs[i+1], "")
	}
	return cli.NewContext(cli.NewApp(), fs, nil)
}

func tcpEchoServer() func() {
	listener, err := net.Listen("tcp", sshdAddr)
	noerr(err)

	go func() {
		for {
			conn, err := listener.Accept()
			// Because historically they have not exported the error that they return. See issues #4373 and #19252.
			if e, ok := err.(*net.OpError); ok && e.Unwrap().Error() == "use of closed network connection" { //nolint:errorlint
				break
			}
			noerr(err)
			// naive approach without any goroutines
			reader := bufio.NewReader(conn)
			for {
				message, err := reader.ReadString('\n')
				if err == io.EOF { //nolint:errorlint
					break
				}
				noerr(err)
				// fmt.Println("ECHO SERVER >>>", string(message)) // debug
				_, err = conn.Write([]byte("echo: " + message + "\n"))
				noerr(err)
			}
		}
	}()

	return func() {
		noerr(listener.Close())
	}
}

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}
