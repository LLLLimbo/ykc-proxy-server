package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ACCEPT_MIN_SLEEP = 10 * time.Millisecond
	// ACCEPT_MAX_SLEEP is the maximum acceptable sleep times on temporary errors
	ACCEPT_MAX_SLEEP = 1 * time.Second
)

type Options struct {
	Host                         string
	TcpPort                      int
	HttpPort                     int
	AutoVerification             bool
	AutoHeartbeatResponse        bool
	AutoBillingModelVerify       bool
	AutoTransactionRecordConfirm bool
	MessagingServerType          string
	Servers                      []string
	Username                     string
	Password                     string
	MessageForwarder             MessageForwarder
	PublishSubjectPrefix         string
}

type Server struct {
	Opt       *Options
	Forwarder *MessageForwarder
	Running   bool
	Mu        sync.RWMutex
	QuitCh    chan struct{}
	GrMu      sync.Mutex
	GrRunning bool
	GrWG      sync.WaitGroup
	Done      chan bool
	Shutdown  bool
}

func NewServer(opts *Options) (*Server, error) {
	s := &Server{
		Opt: opts,
	}
	return s, nil
}

func (s *Server) Start() {
	o := s.Opt

	var hl net.Listener
	var err error

	port := o.TcpPort
	if port == -1 {
		port = 0
	}
	hp := net.JoinHostPort(o.Host, strconv.Itoa(port))
	s.Mu.Lock()
	if s.Shutdown {
		s.Mu.Unlock()
		return
	}
	hl, err = net.Listen("tcp", hp)
	if err != nil {
		s.Mu.Unlock()
		log.Fatalf("Unable to listen for tcp connections: %v", err)
		return
	}
	if port == 0 {
		o.TcpPort = hl.Addr().(*net.TCPAddr).Port
	}

	//go s.acceptConnections(hl, "YKC", func(conn net.Conn) { s.createClient(conn, nil) }, nil)
	s.Mu.Unlock()
}

// Protected check on running state
func (s *Server) isRunning() bool {
	s.Mu.RLock()
	running := s.Running
	s.Mu.RUnlock()
	return running
}

// The following code is modified from nats https://github.com/nats-io/nats-server.git
func (s *Server) acceptConnections(l net.Listener, acceptName string, createFunc func(conn net.Conn), errFunc func(err error) bool) {
	tmpDelay := ACCEPT_MIN_SLEEP

	for {
		conn, err := l.Accept()
		if err != nil {
			if errFunc != nil && errFunc(err) {
				return
			}
			if tmpDelay = s.acceptError(acceptName, err, tmpDelay); tmpDelay < 0 {
				break
			}
			continue
		}
		tmpDelay = ACCEPT_MIN_SLEEP
		if !s.startGoRoutine(func() {
			createFunc(conn)
			s.GrWG.Done()
		}) {
			conn.Close()
		}
	}
	log.Debugf(acceptName + " accept loop exiting..")
	s.Done <- true
}

// If given error is a net.Error and is temporary, sleeps for the given
// delay and double it, but cap it to ACCEPT_MAX_SLEEP. The sleep is
// interrupted if the server is shutdown.
// An error message is displayed depending on the type of error.
// Returns the new (or unchanged) delay, or a negative value if the
// server has been or is being shutdown.
func (s *Server) acceptError(acceptName string, err error, tmpDelay time.Duration) time.Duration {
	if !s.isRunning() {
		return -1
	}
	//lint:ignore SA1019 We want to retry on a bunch of errors here.
	if ne, ok := err.(net.Error); ok && ne.Temporary() { // nolint:staticcheck
		log.Errorf("Temporary %s Accept Error(%v), sleeping %dms", acceptName, ne, tmpDelay/time.Millisecond)
		select {
		case <-time.After(tmpDelay):
		case <-s.QuitCh:
			return -1
		}
		tmpDelay *= 2
		if tmpDelay > ACCEPT_MAX_SLEEP {
			tmpDelay = ACCEPT_MAX_SLEEP
		}
	} else {
		log.Errorf("%s Accept error: %v", acceptName, err)
	}
	return tmpDelay
}

func (s *Server) startGoRoutine(f func()) bool {
	var started bool
	s.GrMu.Lock()
	if s.GrRunning {
		s.GrWG.Add(1)
		go f()
		started = true
	}
	s.GrMu.Unlock()
	return started
}

func parseOptions() *Options {
	host := flag.String("host", "0.0.0.0", "host")
	tcpPort := flag.Int("tcpPort", 27600, "tcpPort")
	httpPort := flag.Int("httpPort", 9556, "httpPort")
	autoVerification := flag.Bool("autoVerification", false, "autoVerification")
	autoHeartbeatResponse := flag.Bool("autoHeartbeatResponse", true, "autoHeartbeatResponse")
	autoBillingModelVerify := flag.Bool("autoBillingModelVerify", false, "autoBillingModelVerify")
	autoTransactionRecordConfirm := flag.Bool("autoTransactionRecordConfirm", false, "autoTransactionRecordConfirm")
	messagingServerType := flag.String("messagingServerType", "", "messagingServerType")
	servers := flag.String("servers", "", "servers")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	flag.Parse()

	//splitting servers with comma
	serversArr := strings.Split(*servers, ",")

	opt := &Options{
		Host:                         *host,
		TcpPort:                      *tcpPort,
		HttpPort:                     *httpPort,
		AutoVerification:             *autoVerification,
		AutoHeartbeatResponse:        *autoHeartbeatResponse,
		AutoBillingModelVerify:       *autoBillingModelVerify,
		AutoTransactionRecordConfirm: *autoTransactionRecordConfirm,
		MessagingServerType:          *messagingServerType,
		Servers:                      serversArr,
		Username:                     *username,
		Password:                     *password,
	}
	return opt
}
