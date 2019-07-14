package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"os/exec"

	quic "github.com/lucas-clemente/quic-go"
	"golang.org/x/net/context"
	cli "gopkg.in/urfave/cli.v2"
)

func server(c *cli.Context) error {
	// generate TLS certificate
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quicssh"},
	}

	// configure listener
	listener, err := quic.ListenAddr(c.String("bind"), config, nil)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Listening at %q...", c.String("bind"))

	ctx := context.Background()
	for {
		log.Printf("Accepting connection...")
		session, err := listener.Accept(ctx)
		if err != nil {
			log.Printf("listener error: %v", err)
			continue
		}

		go serverSessionHandler(ctx, session)
	}
	return nil
}

func serverSessionHandler(ctx context.Context, session quic.Session) {
	log.Printf("hanling session...")
	defer session.Close()
	for {
		stream, err := session.AcceptStream(ctx)
		if err != nil {
			log.Printf("session error: %v", err)
			break
		}
		go serverStreamHandler(ctx, stream)
	}
}

func serverStreamHandler(ctx context.Context, conn io.ReadWriteCloser) {
	fmt.Fprintf(conn, "hello, welcome!\n")
	log.Printf("handling stream...")
	defer conn.Close()
	cmd := exec.Command("rev")
	cmd.Stdin = conn
	/*stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}*/
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	readAndWrite(ctx, conn, cmd.Stdout, nil)
	//go netCopy(conn, cmd.Stdout)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("FINISHED")
}

func netCopy(input io.Reader, output io.Writer) (err error) {
	buf := make([]byte, 8192)
	for {
		count, err := input.Read(buf)
		if err != nil {
			if err == io.EOF && count > 0 {
				output.Write(buf[:count])
			}
			break
		}
		if count > 0 {
			fmt.Println(buf, count)
			output.Write(buf[:count])
		}
	}
	return
}

/*
log.Printf("Spawning subcommand...")
		//cmd := exec.Command("rev")
		//cmd.Stdout = stream
		//cmd.Stdin = stream
		//cmd.Run()

		for {
			log.Printf("Accepting stream...")
			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				return err
			}
			fmt.Println(stream)
		}
	}
	return nil
*/
//ctx, cancel := context.WithCancel(context.Background())
/*
		stdout, err := spawn.StdoutPipe()
		if err != nil {
			return err
		}
		stdin, err := spawn.StdinPipe()
		if err != nil {
			return err
		}
		var wg sync.WaitGroup
		wg.Add(2)
	log.Printf("Piping stdin+stdout with strean...")
	c1 := readAndWrite(ctx, stream, stdin, &wg)
	c2 := readAndWrite(ctx, stdout, stream, &wg)
	if err := spawn.Start(); err != nil {
		return err
	}
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
	fmt.Fprintf(stream, "BEFORE")
*/
/*fmt.Fprintf(stream, "BEFORE")
if _, err = io.Copy(stream, stream); err != nil {
	return err
}
fmt.Fprintf(stream, "AFTER")*/
//return spawn.Wait()
