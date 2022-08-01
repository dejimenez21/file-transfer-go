package main

import (
	"fmt"
	"log"
)

type channel struct {
	name             string
	suscribedClients []client
}

func (c *channel) addClient(newClient client) {
	c.suscribedClients = append(c.suscribedClients, newClient)
}

func (c *channel) broadcast(cmd command) {
	clients := make([]client, len(c.suscribedClients))
	copy(clients, c.suscribedClients)
	cftpBytes, err := serializeDelivery(cmd)
	if err != nil {
		err = fmt.Errorf("error serializing %s delivery throug %s for : %v", cmd.FileInfo.Name, c.name, err)
		log.Println(err)
		return
	}
	for _, client := range clients {
		go c.deliver(cftpBytes, client)
	}
}

func (c *channel) deliver(cftpBytes []byte, receiver client) {
	receiver.writeDelivery(cftpBytes)
}
