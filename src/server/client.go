package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"server/cftp"
	"server/cftp/models"
	"strings"
)

type Client struct {
	Conn        net.Conn
	CmdChan     chan models.Request
	ContentChan chan []byte
	Disconnect  chan *Client
}

func newClient(conn net.Conn) *Client {
	return &Client{
		Conn:        conn,
		CmdChan:     make(chan models.Request),
		ContentChan: make(chan []byte),
		Disconnect:  make(chan *Client),
	}
}

func (c *Client) ReadRequest() {
	for {
		reader := bufio.NewReader(c.Conn)
		data, err := reader.ReadBytes(models.END_OF_MSG)
		if err != nil {
			log.Printf("Client %v disconected", c.Conn.RemoteAddr())
			c.Disconnect <- c
			return
		}
		stringCmd := string(data)
		cmd, err := cftp.DeserializeRequest(strings.TrimSuffix(stringCmd, string(models.END_OF_MSG)))
		if err != nil {
			log.Printf("Error deserializing message: %v", err)
		}
		c.CmdChan <- cmd

		if cmd.Method == models.REQ_SEND {
			for i := 0; i < int(cmd.FileInfo.Size); {
				contentData, err := c.readFileContent(DEFAULT_BUFFER_SIZE, reader)
				if err != nil {
					log.Printf("Client %v disconected", c.Conn.RemoteAddr())
					c.Disconnect <- c
					return
				}
				i += len(contentData)
				c.ContentChan <- contentData
			}
		}
	}
}

func (c *Client) Write(bytes []byte) error {
	_, err := c.Conn.Write(bytes)
	if err != nil {
		c.Disconnect <- c
		return fmt.Errorf("Client %v disconected: %v", c.Conn.RemoteAddr(), err)
	}
	return nil
}

func (c *Client) readFileContent(bufSize int, reader *bufio.Reader) (data []byte, err error) {
	data = make([]byte, bufSize)
	n, err := reader.Read(data)
	data = data[:n]
	return
}
