package main

import (
	"bufio"
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
	WriteChan   chan []byte
	Disconnect  chan *Client
}

func newClient(conn net.Conn) *Client {
	return &Client{
		Conn:        conn,
		CmdChan:     make(chan models.Request),
		ContentChan: make(chan []byte),
		WriteChan:   make(chan []byte),
		Disconnect:  make(chan *Client),
	}
}

func (c *Client) ReadRequest() {
	for {
		reader := bufio.NewReader(c.Conn)
		data, err := reader.ReadBytes(cftp.END_OF_MSG)
		if err != nil {
			log.Printf("Client %v disconected", c.Conn.RemoteAddr())
			c.Disconnect <- c
			return
		}
		stringCmd := string(data)
		cmd, err := cftp.DeserializeCommand(strings.TrimSuffix(stringCmd, string(cftp.END_OF_MSG)))
		if err != nil {
			log.Printf("Error deserializing message: %v", err)
		}
		c.CmdChan <- cmd

		if cmd.Method == CMD_SEND {
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

func (c *Client) StartWriter() {
	for {
		msg := <-c.WriteChan
		_, err := c.Conn.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
	}
	c.Disconnect <- c
}

func (c *Client) readFileContent(bufSize int, reader *bufio.Reader) (data []byte, err error) {
	data = make([]byte, bufSize)
	n, err := reader.Read(data)
	data = data[:n]
	return
}
