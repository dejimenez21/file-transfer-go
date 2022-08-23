package main

import (
	"fmt"
	"log"
	"server/cftp"
	"server/cftp/models"
	"sync"
)

type channel struct {
	name                 string
	suscribedClients     map[string]*client
	suscribedClientsLock sync.RWMutex
}

func (c *channel) addClient(newClient *client) {
	c.suscribedClientsLock.Lock()
	c.suscribedClients[newClient.conn.RemoteAddr().String()] = newClient
	c.suscribedClientsLock.Unlock()
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
		client.writeChan <- cftpBytes
	}
	var chunkSeq int64 = 0
	deliveryID := models.NewDeliveryId()
	for {
		fileContent := <-contentChan
		chunkSeq++
		del := models.Delivery{Content: fileContent, ID: deliveryID, Seq: chunkSeq, Size: len(fileContent)}
		deliveryBytes := cftp.SerializeChunkDelivery(del)
		for _, client := range clients {
			client.writeChan <- deliveryBytes
		}
	}

}

func (c *channel) copySuscribedClients() map[string]*client {
	copy := make(map[string]*client)
	c.suscribedClientsLock.RLock()
	for k, v := range c.suscribedClients {
		copy[k] = v
	}
	c.suscribedClientsLock.RUnlock()
	return copy
}

func (c *channel) UnsuscribeClient(client *client) {
	c.suscribedClientsLock.Lock()
	delete(c.suscribedClients, client.conn.RemoteAddr().String())
	c.suscribedClientsLock.Unlock()
}
