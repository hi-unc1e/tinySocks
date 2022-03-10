package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type ProxyAuth struct {
	Enable   bool
	Username string
	Passwd   string
}

func DialTcpByProxy(proxyStr string, addr string) (c net.Conn, err error) {
	var proxyUrl *url.URL
	if proxyUrl, err = url.Parse(proxyStr); err != nil {
		return
	}

	auth := &ProxyAuth{}
	if proxyUrl.User != nil {
		auth.Enable = true
		auth.Username = proxyUrl.User.Username()
		auth.Passwd, _ = proxyUrl.User.Password()
	}

	switch proxyUrl.Scheme {
	case "http":
		return DialTcpByHttpProxy(proxyUrl.Host, addr, auth)
	case "socks5":
		return DialTcpBySocks5Proxy(proxyUrl.Host, addr, auth)
	default:
		err = fmt.Errorf("Proxy URL scheme must be http or socks5, not [%s]", proxyUrl.Scheme)
		return
	}
}

func DialTcpByHttpProxy(proxyHost string, dstAddr string, auth *ProxyAuth) (c net.Conn, err error) {
	if c, err = net.Dial("tcp", proxyHost); err != nil {
		return
	}

	req, err := http.NewRequest("CONNECT", "http://"+dstAddr, nil)
	if err != nil {
		return
	}
	if auth.Enable {
		req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth.Username+":"+auth.Passwd)))
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Write(c)

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("DialTcpByHttpProxy error, StatusCode [%d]", resp.StatusCode)
		return
	}
	return
}

func DialTcpBySocks5Proxy(proxyHost string, dstAddr string, auth *ProxyAuth) (c net.Conn, err error) {
	var s5Auth *proxy.Auth
	if auth.Enable {
		s5Auth = &proxy.Auth{
			User:     auth.Username,
			Password: auth.Passwd,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", proxyHost, s5Auth, nil)
	if err != nil {
		return nil, err
	}

	if c, err = dialer.Dial("tcp", dstAddr); err != nil {
		return
	}
	return
}

func main() {
	relayServer := flag.String("s", "", "The ligolo server ip:port (e.g. example.com:443)")
	targetServer := flag.String("t", "", "The destination server ip:port (e.g. 192.168.1.3:3389, 192.168.1.3:22, etc.) - when not specified, Ligolo starts a socks5 proxy server")
	proxyStr := flag.String("proxy", "", "Use proxy to connect ligolo server(e.g. http://user:passwd@192.168.1.128:8080 socks5://user:passwd@192.168.1.128:1080)")
	flag.Parse()
	for {
		err := StartLigolo(*relayServer, *targetServer, *proxyStr)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Warning("Restarting ligolo client...")
		time.Sleep(10 * time.Second)
	}
}

func StartLigolo(relayServer string, targetServer string, proxyStr string) error {
	var socks *socks5.Server
	logrus.Infoln("Connecting to ligolo server...")
	config := &tls.Config{InsecureSkipVerify: true}
	var conn net.Conn
	var err error

	if proxyStr == "" {
		conn, err = tls.Dial("tcp", relayServer, config)
		if err != nil {
			return err
		}

	} else {
		logrus.Infoln("Using proxy")
		conn, err = DialTcpByProxy(proxyStr, relayServer)
		if err != nil {
			return err
		}
		conn = tls.Client(conn, config)
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
		if targetServer == "" {
			go socks.ServeConn(stream)
		} else {
			proxyConn, err := net.Dial("tcp", targetServer)
			if err != nil {
				logrus.Errorf("Error creating Proxy TCP connection ! Error : %s\n", err)
				return err
			}
			go handleRelay(stream, proxyConn)
		}
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
