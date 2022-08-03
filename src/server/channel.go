package main

import (
	"fmt"
	"log"
)

type channel struct {
	name             string
	suscribedClients map[string]*client
}

func (c *channel) addClient(newClient *client) {
	// c.suscribedClients = append(c.suscribedClients, newClient)
	//TODO: Add a lock to the suscribedClients map
	c.suscribedClients[newClient.conn.RemoteAddr().String()] = newClient
}

func (c *channel) broadcast(cmd command, contentChan chan []byte) {
	log.Printf("Broadcasting file from %s through %s", cmd.Meta.SenderAddress, c.name)
	clients := c.copySuscribedClients()
	cftpBytes, err := serializeCommand(cmd)
	if err != nil {
		err = fmt.Errorf("error serializing %s delivery throug %s for : %v", cmd.FileInfo.Name, c.name, err)
		log.Println(err)
		return
	}
	for _, client := range clients {
		// go c.deliver(cftpBytes, client)
		client.writeChan <- cftpBytes
	}
	var chunkSeq int64 = 0
	for {
		fileContent := <-contentChan
		chunkSeq++
		del := delivery{Content: fileContent, DeliveryId: c.newDeliveryId(), Seq: chunkSeq, Size: len(fileContent)}
		delBytes := serializeChunkDelivery(del)
		for _, client := range clients {
			client.writeChan <- delBytes
		}
	}

}

func (c *channel) deliver(cftpBytes []byte, receiver *client) {
	// err := (*receiver).writeDelivery(cftpBytes)
	// if err != nil {
	// 	log.Printf("could not deliver to client %s: %v", receiver.conn.RemoteAddr().String(), err)
	// 	//TODO: add a lock
	// 	delete(c.suscribedClients, receiver.conn.RemoteAddr().String())
	// }
	log.Printf("Message delivered to %s", (*receiver).conn.RemoteAddr().String())
}

func (c *channel) copySuscribedClients() map[string]*client {
	copy := make(map[string]*client)
	for k, v := range c.suscribedClients {
		copy[k] = v
	}
	return copy
}

func (c *channel) newDeliveryId() int64 {
	//TODO: implement struct member to keep track of ids
	return 1
}
