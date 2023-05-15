package main

import (
	"errors"
	"net"
	"sync"
)

var clients sync.Map

func StoreClient(id string, conn net.Conn) {
	clients.Store(id, conn)
}

func GetClient(id string) (net.Conn, error) {
	value, ok := clients.Load(id)
	if ok {
		conn := value.(net.Conn)
		return conn, nil
	} else {
		return nil, errors.New("client does not exist")
	}
}
