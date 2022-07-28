package main

import (
	"bufio"
	"log"
	"net"
)

type client struct {
	conn    net.Conn
	cmdChan chan command
}

func (c *client) readCommand(cmdChn chan command) {
	for {
		data, err := bufio.NewReader(c.conn).ReadBytes(0x04)
		if err != nil {
			log.Printf("Error reading message from: %v", c.conn.RemoteAddr())
			return
		}

		cmd, err := deserializeCommand(string(data))
		if err != nil {
			log.Printf("Error deserializing message: %v", err)
		}

		cmdChn <- cmd
	}
}

func (c *client) writeDelivery(message []byte) {
	c.conn.Write(message)
}
