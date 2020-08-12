package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

var tlsFingerprint string

var (
	ErrInvalidServerCert = fmt.Errorf("invalid TLS server certificate")
	ErrInvalidPinnedCert = fmt.Errorf("invalid TLS pinned certificate")
)

func main() {
	relayServer := flag.String("s", "example.com:443", "The ligolo server (the connect-back address)")
	flag.Parse()
	for {
		err := StartLigolo(*relayServer)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Warning("Restarting ligolo client...")
		time.Sleep(10 * time.Second)
	}
}

func StartLigolo(relayServer string) error {
	var socks *socks5.Server
	logrus.Infoln("Connecting to ligolo server...")

	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", relayServer, config)
	if err != nil {
		return err
	}
	socks, err = startSocksProxy()
	if err != nil {
		logrus.Error("Could not start SOCKS5 proxy !")
		return err
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}

	logrus.Infoln("Waiting for connections....")

	for {
		stream, err := session.Accept()
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{"active_sessions": session.NumStreams()}).Println("Accepted new connection !")
		// When no targetServer are specified, starts a socks5 proxy
		go socks.ServeConn(stream)
	}
}

func startSocksProxy() (*socks5.Server, error) {
	conf := &socks5.Config{}
	socks, err := socks5.New(conf)
	if err != nil {
		logrus.Error("Could not start SOCKS5 proxy !")
		return nil, err
	}
	return socks, nil
}

func handleRelay(src net.Conn, dst net.Conn) {
	stop := make(chan bool, 2)

	go relay(src, dst, stop)
	go relay(dst, src, stop)

	select {
	case <-stop:
		return
	}
}

func relay(src net.Conn, dst net.Conn, stop chan bool) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
	stop <- true
	return
}
