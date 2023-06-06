package main

import (
	"github.com/go-resty/resty/v2"
	"math/rand"
	"strings"
	"time"
)
import log "github.com/sirupsen/logrus"

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

	client := resty.New()
	var res string
	client.SetTimeout(3 * time.Second)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&res).
		SetBody(message).
		Post(e + "/" + mid)
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
	// ...
}

func (h *NatsForwarder) Subscribe(topic string, handler func(message []byte)) error {
	//TODO implement me
	panic("implement me")
}

func (h *NatsForwarder) Close() error {
	//TODO implement me
	panic("implement me")
}

func (h *NatsForwarder) Connect() error {
	//...
	return nil
}

func (h *NatsForwarder) Publish(mid string, message []byte) error {
	return nil
}
