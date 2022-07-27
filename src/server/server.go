package main

import (
	"log"
	"net"
)

type server struct {
	channels         []channel
	connectedClients []client
}

func (s *server) newClient(conn net.Conn) {
	log.Println("Client connected from", conn.RemoteAddr())
	newClient := &client{conn: conn}
	newClient.readCommand()
}
