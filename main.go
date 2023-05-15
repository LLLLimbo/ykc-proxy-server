package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type ProxyEnv struct {
	Host string
	Port int
}

func main() {
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
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Error("error reading ", err)
		return err
	}

	hex := BytesToHex(buffer[:n])

	//is encrypted ?
	encrypted := false
	if hex[4] == "1" {
		encrypted = true
	}

	//message length
	length, err := strconv.ParseInt(hex[1], 16, 64)

	//message sequence number
	seq, err := strconv.ParseInt(hex[3]+hex[2], 16, 64)

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
		"frame_id":  hex[5],
	}).Info("received message")

	switch hex[5] {
	case "01":
		VerificationRouter(hex, header, conn)
		break
	case "03":
		HeartbeatRouter(hex, header, conn)
		break
	case "05":
		BillingModelVerificationRouter(hex, header, conn)
		break
	case "13":
		OfflineDataReportMessageRouter(hex, header)
	case "33":
		RemoteBootstrapResponseRouter(hex, header)
	case "35":
		RemoteShutdownResponseRouter(hex, header)
	case "3b":
		TransactionRecordMessageRouter(buffer, hex, header)
	default:
		break

	}
	return nil
}
