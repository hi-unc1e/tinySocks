package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"io"
	Ligolo "ligolo"
	"net"
)

func main() {
	localServer := flag.String("p", "127.0.0.1:1080", "The local socks5 server address or ip:port to connect")
	relayServer := flag.String("l", "0.0.0.0:443", "The relay server listening address (the connect-back address)")
	certFile := flag.String("cert", "cert.pem", "The TLS server certificate,Unnecessary")
	keyFile := flag.String("key", "key.pem", "The TLS server key,Unnecessary")

	flag.Parse()

	relay := NewLigoloRelay(*localServer, *relayServer, *certFile, *keyFile)
	relay.Start()
}

// LigoloRelay structure contains configuration, the current session and the ConnectionPool
type LigoloRelay struct {
	LocalServer    string
	RelayServer    string
	CertFile       string
	KeyFile        string
	ConnectionPool chan *yamux.Session
	Session        *yamux.Session
}

// NewLigoloRelay creates a new LigoloRelay struct
func NewLigoloRelay(localServer string, relayServer string, certFile string, keyFile string) *LigoloRelay {
	return &LigoloRelay{LocalServer: localServer, RelayServer: relayServer, CertFile: certFile, KeyFile: keyFile, ConnectionPool: make(chan *yamux.Session, 100)}
}

// Start listening for local and relay connections
func (ligolo LigoloRelay) Start() {

	logrus.WithFields(logrus.Fields{"socks5server": ligolo.LocalServer, "ligoloserver": ligolo.RelayServer}).Println("Ligolo server started.")
	go ligolo.startRelayHandler()
	ligolo.startLocalHandler()
}

// Listen for Ligolo connections
func (ligolo LigoloRelay) startRelayHandler() {

	cer, err := tls.LoadX509KeyPair(ligolo.CertFile, ligolo.KeyFile)
	if err != nil {
		cer, _ = tls.X509KeyPair([]byte(Ligolo.CertPEM), []byte(Ligolo.KeyPEM))
		//
		//logrus.Warning("Could not load TLS certificate.")
		//return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp4", ligolo.RelayServer, config)
	if err != nil {
		logrus.Errorf("Could not bind to port : %v\n", err)
		return
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			logrus.Errorf("Could not accept connection : %v\n", err)
			return
		}

		session, err := handleRelayConnection(c)
		if err != nil {
			logrus.Errorf("Could not start session : %v\n", err)
			continue
		}
		ligolo.ConnectionPool <- session
	}

}

// Listen for local connections
func (ligolo LigoloRelay) startLocalHandler() {
	listener, err := net.Listen("tcp4", ligolo.LocalServer)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	ligolo.Session = <-ligolo.ConnectionPool
	go func() {
		for {
			<-ligolo.Session.CloseChan()
			logrus.WithFields(logrus.Fields{"remoteaddr": ligolo.Session.RemoteAddr()}).Println("Received session shutdown.")
			ligolo.Session = <-ligolo.ConnectionPool
			logrus.WithFields(logrus.Fields{"remoteaddr": ligolo.Session.RemoteAddr()}).Println("New session acquired.")
		}
	}()
	logrus.Println("Session acquired. Starting relay.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go ligolo.handleLocalConnection(conn)
	}
}

// Handle new local connections
func (ligolo LigoloRelay) handleLocalConnection(conn net.Conn) {
	if ligolo.Session.IsClosed() {
		logrus.Warning("Closing connection because no session available !")
		conn.Close()
		return
	}

	logrus.Println("New proxy connection. Establishing new session.")

	stream, err := ligolo.Session.Open()
	if err != nil {
		logrus.Errorf("Could not open session : %s\n", err)
		return
	}

	logrus.Println("Yamux session established.")

	go relay(conn, stream)
	go relay(stream, conn)

}

// Handle new ligoloc connections
func handleRelayConnection(conn net.Conn) (*yamux.Session, error) {
	logrus.WithFields(logrus.Fields{"remoteaddr": conn.RemoteAddr().String()}).Info("New relay connection.\n")
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return nil, err
	}
	ping, err := session.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Printf("Session ping : %v\n", ping)
	return session, nil
}

func relay(src net.Conn, dst net.Conn) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
	return
}
