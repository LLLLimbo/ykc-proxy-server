package main

import "github.com/go-resty/resty/v2"
import log "github.com/sirupsen/logrus"

type MessageForwarder interface {
	Connect() error
	Publish(topic string, message []byte) error
	Subscribe(topic string, handler func(message []byte)) error
	Close() error
}

type HTTPForwarder struct {
	// ...
}

func (h *HTTPForwarder) Publish(topic string, message []byte) error {
	client := resty.New()
	var res string
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&res).
		SetBody(message).
		Post(topic)
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
