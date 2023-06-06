/*
The TCP Server implementation is referenced from nats
https://github.com/nats-io/nats-server.git
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	ACCEPT_MIN_SLEEP = 10 * time.Millisecond
	// ACCEPT_MAX_SLEEP is the maximum acceptable sleep times on temporary errors
	ACCEPT_MAX_SLEEP = 1 * time.Second
)

type Options struct {
	Host                   string
	Port                   int
	AutoVerification       bool
	AutoHeartbeatResponse  bool
	AutoBillingModelVerify bool
	MessagingServerType    string
	Servers                []string
	Username               string
	Password               string
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

	port := o.Port
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
		log.Fatalf("unable to listen for tcp connections: %v", err)
		return
	}
	if port == 0 {
		o.Port = hl.Addr().(*net.TCPAddr).Port
	}

	go s.acceptConnections(hl, "ykc", nil)
	s.Mu.Unlock()
}

// Protected check on running state
func (s *Server) isRunning() bool {
	s.Mu.RLock()
	running := s.Running
	s.Mu.RUnlock()
	return running
}

func (s *Server) acceptConnections(l net.Listener, acceptName string, errFunc func(err error) bool) {
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
			createFunc(s, conn)
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

func createFunc(s *Server, conn net.Conn) {
	defer conn.Close()

	log.WithFields(log.Fields{
		"address": conn.RemoteAddr().String(),
	}).Info("new client connected")

	var connErr error
	for connErr == nil {
		connErr = drain_(s, conn)
		time.Sleep(time.Millisecond * 1)
	}
}

func drain_(s *Server, conn net.Conn) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("error reading ", err)
		return err
	}

	hex := BytesToHex(buf[:n])

	//is encrypted ?
	encrypted := false
	if buf[4] == byte(0x01) {
		encrypted = true
	}

	//message length
	length := buf[1]

	//message sequence number
	seq := buf[3]<<8 | buf[2]

	header := &Header{
		Length:    int(length),
		Seq:       int(seq),
		Encrypted: encrypted,
		FrameId:   hex[5],
	}

	log.WithFields(log.Fields{
		"hex":       hex,
		"encrypted": encrypted,
		"length":    length,
		"seq":       seq,
		"frame_id":  int(buf[5]),
	}).Info("received message")

	switch buf[5] {
	case Verification:
		VerificationRouter(buf, hex, header, conn)
		break
	case Heartbeat:
		HeartbeatRouter(hex, header, conn)
		break
	case BillingModelVerification:
		BillingModelVerificationRouter(hex, header, conn)
		break
	case OfflineDataReport:
		OfflineDataReportMessageRouter(hex, header)
	case RemoteBootstrapResponse:
		RemoteBootstrapResponseRouter(hex, header)
	case RemoteShutdownResponse:
		RemoteShutdownResponseRouter(hex, header)
	case TransactionRecord:
		TransactionRecordMessageRouter(buf, hex, header)
	default:
		break

	}
	return nil
}
