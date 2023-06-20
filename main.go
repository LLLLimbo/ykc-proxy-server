package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	opt := parseOptions()

	//define message forwarder
	var f MessageForwarder
	switch opt.MessagingServerType {
	case "http":
		f := &HTTPForwarder{
			Endpoints: opt.Servers,
		}
		opt.MessageForwarder = f
		break
	case "nats":
		servers := strings.Join(opt.Servers, ",")
		f := &NatsForwarder{
			Servers:  servers,
			Username: opt.Username,
			Password: opt.Password,
		}
		opt.MessageForwarder = f
		break
	default:
		f = nil
		opt.MessageForwarder = f
		break
	}

	go enableTcpServer(opt)
	go enableHttpServer(opt)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	log.Info("exit:", sig)
	os.Exit(0)
}

func enableTcpServer(opt *Options) {
	host := opt.Host
	port := strconv.Itoa(opt.TcpPort)
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		log.Error("error resolving address:", err)
		return
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("error listening:", err)
		return
	}
	defer ln.Close()
	log.Info("server listening on", addr.String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error("error accepting connection:", err)
			continue
		}
		StoreClient(conn.RemoteAddr().String(), conn)
		go handleConnection(opt, conn)
	}
}

func enableHttpServer(opt *Options) {
	r := gin.Default()
	r.POST("/proxy/02", VerificationResponseRouter)
	r.POST("/proxy/06", BillingModelVerificationResponseRouter)
	r.POST("/proxy/34", RemoteBootstrapRequestRouter)
	r.POST("/proxy/36", RemoteShutdownRequestRouter)
	r.POST("/proxy/40", TransactionRecordConfirmedRouter)
	r.POST("/proxy/58", SetBillingModelRequestRouter)
	r.POST("/proxy/92", RemoteRebootRequestMessageRouter)

	host := opt.Host
	port := strconv.Itoa(opt.HttpPort)
	err := r.Run(host + ":" + port)
	if err != nil {
		panic(err)
	}
}

func handleConnection(opt *Options, conn net.Conn) {
	defer conn.Close()

	log.WithFields(log.Fields{
		"address": conn.RemoteAddr().String(),
	}).Info("new client connected")

	var connErr error
	for connErr == nil {
		connErr = drain(opt, conn)
		time.Sleep(time.Millisecond * 1)
	}

}

func drain(opt *Options, conn net.Conn) error {
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
		VerificationRouter(opt, buf, hex, header, conn)
		break
	case Heartbeat:
		HeartbeatRouter(opt, hex, header, conn)
		break
	case BillingModelVerification:
		BillingModelVerificationRouter(opt, hex, header, conn)
		break
	case OfflineDataReport:
		OfflineDataReportMessageRouter(opt, hex, header)
		break
	case RemoteBootstrapResponse:
		RemoteBootstrapResponseRouter(opt, hex, header)
		break
	case RemoteShutdownResponse:
		RemoteShutdownResponseRouter(opt, hex, header)
		break
	case SetBillingModelResponse:
		SetBillingModelResponseMessageRouter(opt, hex, header)
		break
	case RemoteRebootResponse:
		RemoteRebootResponseMessageRouter(opt, hex, header)
		break
	case TransactionRecord:
		TransactionRecordMessageRouter(opt, buf, hex, header)
		break
	default:
		break

	}
	return nil
}
