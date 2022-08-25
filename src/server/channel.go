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
	suscribedClients     map[string]*Client
	suscribedClientsLock sync.RWMutex
}

func (c *channel) AddClient(newClient *Client) {
	c.suscribedClientsLock.Lock()
	c.suscribedClients[newClient.Conn.RemoteAddr().String()] = newClient
	c.suscribedClientsLock.Unlock()
}

func (c *channel) Broadcast(cmd models.Request, contentChan chan []byte) {
	log.Printf("Broadcasting file from %s through %s", cmd.Meta.SenderAddress, c.name)
	clients := c.copySuscribedClients()
	cftpBytes, err := cftp.SerializeCommand(cmd)
	if err != nil {
		err = fmt.Errorf("error serializing %s delivery throug %s for : %v", cmd.FileInfo.Name, c.name, err)
		log.Println(err)
		return
	}
	for _, client := range clients {
		client.Write(cftpBytes)
	}
	var chunkSeq int64 = 0
	deliveryID := cmd.Meta.RequestId
	for {
		fileContent := <-contentChan
		chunkSeq++
		del := models.Delivery{Content: fileContent, ID: deliveryID, Seq: chunkSeq, Size: len(fileContent)}
		deliveryBytes := cftp.SerializeChunkDelivery(del)
		for _, client := range clients {
			err := client.Write(deliveryBytes)
			if err != nil {
				log.Printf("error sending chunk to client: %v", err)
				delete(clients, client.Conn.RemoteAddr().String())
			}
		}
	}

}

func (c *channel) RemoveClient(client *Client) {
	c.suscribedClientsLock.Lock()
	delete(c.suscribedClients, client.Conn.RemoteAddr().String())
	c.suscribedClientsLock.Unlock()
}

func (c *channel) copySuscribedClients() map[string]*Client {
	copy := make(map[string]*Client)
	c.suscribedClientsLock.RLock()
	for k, v := range c.suscribedClients {
		copy[k] = v
	}
	c.suscribedClientsLock.RUnlock()
	return copy
}
