package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go/http3"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name: "server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "bind", Value: "localhost:4242"},
					&cli.StringFlag{Name: "chain", Value: "./test.crt"},
					&cli.StringFlag{Name: "key", Value: "./test.key"},
				},
				Action: server,
			},
			&cli.Command{
				Name: "client",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "addr", Value: "https://localhost:4242"},
				},
				Action: client,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func server(c *cli.Context) error {
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Printf("Listening on %s...", c.String("bind"))
	return http3.ListenAndServeQUIC(c.String("bind"), c.String("chain"), c.String("key"), nil)
}

func client(c *cli.Context) error {
	roundTripper := &http3.RoundTripper{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	defer roundTripper.Close()
	hClient := &http.Client{Transport: roundTripper}
	resp, err := hClient.Get(c.String("addr"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	body := &bytes.Buffer{}
	if _, err = io.Copy(body, resp.Body); err != nil {
		panic(err)
	}
	fmt.Println(body)
	return nil
}
