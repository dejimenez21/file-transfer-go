package main

import (
	"fmt"
	"log"
	"server/cftp"
	"server/cftp/models"
	"sync"
)

type Channel struct {
	Name                 string
	SuscribedClients     map[string]*Client
	suscribedClientsLock sync.RWMutex
}

func newChannel(name string) *Channel {
	return &Channel{
		Name:             name,
		SuscribedClients: make(map[string]*Client),
	}
}

func (c *Channel) AddClient(newClient *Client) {
	c.suscribedClientsLock.Lock()
	c.SuscribedClients[newClient.Conn.RemoteAddr().String()] = newClient
	c.suscribedClientsLock.Unlock()
}

func (c *Channel) Broadcast(cmd models.Request, contentChan chan []byte) {
	log.Printf("Broadcasting file from %s through %s", cmd.Meta.SenderAddress, c.Name)
	clients := c.copySuscribedClients()
	cftpBytes, err := cftp.SerializeRequest(cmd)
	if err != nil {
		err = fmt.Errorf("error serializing %s delivery throug %s for : %v", cmd.FileInfo.Name, c.Name, err)
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

func (c *Channel) RemoveClient(client *Client) {
	c.suscribedClientsLock.Lock()
	delete(c.SuscribedClients, client.Conn.RemoteAddr().String())
	c.suscribedClientsLock.Unlock()
}

func (c *Channel) copySuscribedClients() map[string]*Client {
	copy := make(map[string]*Client)
	c.suscribedClientsLock.RLock()
	for k, v := range c.SuscribedClients {
		copy[k] = v
	}
	c.suscribedClientsLock.RUnlock()
	return copy
}
