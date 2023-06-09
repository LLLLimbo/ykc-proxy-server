package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	"math/rand"
	"strings"
	"time"
)
import log "github.com/sirupsen/logrus"

const (
	NATS_PUBLISH_SUBJECT_PREFIX = "charge.proxy.ykc"
)

type MessageForwarder interface {
	Connect() error
	Publish(mid string, message []byte) error
	Subscribe(topic string, handler func(message []byte)) error
	Close() error
}

type HTTPForwarder struct {
	Endpoints []string
}

func (h *HTTPForwarder) Connect() error {
	panic("implement me")
}

func (h *HTTPForwarder) Subscribe(topic string, handler func(message []byte)) error {
	panic("implement me")
}

func (h *HTTPForwarder) Close() error {
	panic("implement me")
}

func (h *HTTPForwarder) Publish(mid string, message []byte) error {
	//get random endpoint
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(h.Endpoints))
	e := h.Endpoints[randomIndex]
	if strings.HasSuffix(e, "/") {
		e = e[:len(e)-1]
	}

	e = e + "/" + mid

	log.WithFields(log.Fields{
		"mid":      mid,
		"endpoint": e,
	}).Info("forwarding message to http endpoint")

	client := resty.New()
	var res string
	client.SetTimeout(3 * time.Second)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&res).
		SetBody(message).
		Post(e)
	if err != nil {
		log.Errorf("error forwarding message: %v", err)
		return err
	}
	log.WithFields(log.Fields{
		"http_response": res,
	}).Info("response from http forwarder")

	return nil
}

type KafkaForwarder struct {
	// ...
}
type RabbitForwarder struct {
	// ...
}

type NatsForwarder struct {
	Servers  string
	Username string
	Password string
	nc       *nats.Conn
}

func (h *NatsForwarder) Subscribe(topic string, handler func(message []byte)) error {
	//TODO implement me
	panic("implement me")
}

func (h *NatsForwarder) Close() error {
	h.nc.Close()
	return nil
}

func (h *NatsForwarder) Connect() error {
	nc, err := nats.Connect(h.Servers, nats.UserInfo(h.Username, h.Password))
	if err != nil {
		log.Fatalf("can not connect to nats server, error: %s", err.Error())
	}
	h.nc = nc
	return err
}

func (h *NatsForwarder) Publish(mid string, message []byte) error {
	subject := NATS_PUBLISH_SUBJECT_PREFIX + "." + mid
	err := h.nc.Publish(subject, message)
	return err
}
