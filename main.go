package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	parseOptions()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	go enableTcpServer()
	go enableHttpServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	log.Info("exit:", sig)
	os.Exit(0)
}

func parseOptions() *Options {
	var host = flag.String("host", "0.0.0.0", "host")
	var port = flag.Int("port", 27600, "port")
	var autoVerification = flag.Bool("auto-verification", false, "auto verification")
	var autoHeartbeatResponse = flag.Bool("auto-heartbeat-response", false, "auto heartbeat response")
	var autoBillingModelVerify = flag.Bool("auto-billing-model-verify", false, "auto billing model verify")
	var messagingServerType = flag.String("messaging-server-type", "http", "messaging server type,default to http")
	var servers = flag.String("servers", "http://127.0.0.1:8080", "servers,use commas to separate multiple server addresses")
	var username = flag.String("username", "admin", "username")
	var password = flag.String("password", "admin", "password")
	flag.Parse()
	serversArr := strings.Split(*servers, ",")
	options := &Options{
		Host:                   *host,
		Port:                   *port,
		AutoVerification:       *autoVerification,
		AutoHeartbeatResponse:  *autoHeartbeatResponse,
		AutoBillingModelVerify: *autoBillingModelVerify,
		MessagingServerType:    *messagingServerType,
		Servers:                serversArr,
		Username:               *username,
		Password:               *password,
	}
	return options
}

func enableTcpServer() {
	host := "0.0.0.0"
	port := "27600"
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
		go handleConnection(conn)
	}
}

func enableHttpServer() {
	r := gin.Default()
	r.POST("/proxy/02", VerificationResponseRouter)
	r.POST("/proxy/06", BillingModelVerificationResponseRouter)
	r.POST("/proxy/34", RemoteBootstrapRequestRouter)
	r.POST("/proxy/36", RemoteShutdownRequestRouter)
	err := r.Run("0.0.0.0:9556")
	if err != nil {
		panic(err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.WithFields(log.Fields{
		"address": conn.RemoteAddr().String(),
	}).Info("new client connected")

	var connErr error
	for connErr == nil {
		connErr = drain(conn)
		time.Sleep(time.Millisecond * 1)
	}

}

func drain(conn net.Conn) error {
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
