package main

import (
	"fmt"
	"log"
	"server/cftp"
	"server/cftp/models"
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

func (c *channel) broadcast(cmd models.Command, contentChan chan []byte) {
	log.Printf("Broadcasting file from %s through %s", cmd.Meta.SenderAddress, c.name)
	clients := c.copySuscribedClients()
	cftpBytes, err := cftp.SerializeCommand(cmd)
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
	deliveryID := models.NewDeliveryId()
	for {
		fileContent := <-contentChan
		chunkSeq++
		del := models.Delivery{Content: fileContent, ID: deliveryID, Seq: chunkSeq, Size: len(fileContent)}
		delBytes := cftp.SerializeChunkDelivery(del)
		for _, client := range clients {
			client.writeChan <- delBytes
		}
	}

}

func (c *channel) copySuscribedClients() map[string]*client {
	copy := make(map[string]*client)
	for k, v := range c.suscribedClients {
		copy[k] = v
	}
	return copy
}
