# quicssh

> :smile: **quicssh** is a QUIC proxy that allows to use QUIC to connect to an SSH server without needing to patch the client or the server.

[![CircleCI](https://circleci.com/gh/moul/quicssh.svg?style=shield)](https://circleci.com/gh/moul/quicssh)
[![GoDoc](https://godoc.org/moul.io/quicssh?status.svg)](https://godoc.org/moul.io/quicssh)
[![License](https://img.shields.io/github/license/moul/quicssh.svg)](https://github.com/moul/quicssh/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/moul/quicssh.svg)](https://github.com/moul/quicssh/releases)
[![Go Report Card](https://goreportcard.com/badge/moul.io/quicssh)](https://goreportcard.com/report/moul.io/quicssh)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/quicssh.svg)](https://microbadger.com/images/moul/quicssh)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

## Architecture

Standard SSH connection

```
┌───────────────────────────────────────┐             ┌───────────────────────┐
│                  bob                  │             │         wopr          │
│ ┌───────────────────────────────────┐ │             │ ┌───────────────────┐ │
│ │           ssh user@wopr           │─┼────tcp──────┼▶│       sshd        │ │
│ └───────────────────────────────────┘ │             │ └───────────────────┘ │
└───────────────────────────────────────┘             └───────────────────────┘
```

---

SSH Connection proxified with QUIC

```
┌───────────────────────────────────────┐             ┌───────────────────────┐
│                  bob                  │             │         wopr          │
│ ┌───────────────────────────────────┐ │             │ ┌───────────────────┐ │
│ │ssh -o ProxyCommand "quicssh client│ │             │ │       sshd        │ │
│ │     --addr %h:4545" user@wopr     │ │             │ └───────────────────┘ │
│ │                                   │ │             │           ▲           │
│ └───────────────────────────────────┘ │             │           │           │
│                   │                   │             │           │           │
│                process                │             │  tcp to localhost:22  │
│                   │                   │             │           │           │
│                   ▼                   │             │           │           │
│ ┌───────────────────────────────────┐ │             │┌─────────────────────┐│
│ │  quicssh client --addr wopr:4545  │─┼─quic (udp)──▶│   quicssh server    ││
│ └───────────────────────────────────┘ │             │└─────────────────────┘│
└───────────────────────────────────────┘             └───────────────────────┘
```

## Usage

```console
$ quicssh -h
NAME:
   quicssh - A new cli application

USAGE:
   quicssh [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     server
     client
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
   ```

### Client

```console
$ quicssh client -h
NAME:
   quicssh client -

USAGE:
   quicssh client [command options] [arguments...]

OPTIONS:
   --addr value  (default: "localhost:4242")
   --help, -h    show help (default: false)
```

### Server

```console
$ quicssh server -h
NAME:
   quicssh server -

USAGE:
   quicssh server [command options] [arguments...]

OPTIONS:
   --bind value  (default: "localhost:4242")
   --help, -h    show help (default: false)
```

## Install

```console
$ go get -u moul.io/quicssh
```

## License

© 2019 [Manfred Touron](https://manfred.life) -
[Apache-2.0 License](https://github.com/moul/quicssh/blob/master/LICENSE)
