package main

import (
	"bufio"
	"log"
	"net"
	"server/cftp"
	"server/cftp/models"
	"strings"
)

type client struct {
	conn       net.Conn
	writeChan  chan []byte
	disconnect chan<- *client
}

func (c *client) readRequest(cmdChn chan models.Command, contentChn chan<- []byte) (req models.Command) {
	for {
		reader := bufio.NewReader(c.conn)
		data, err := reader.ReadBytes(cftp.END_OF_MSG)
		if err != nil {
			log.Printf("Client %v disconected", c.conn.RemoteAddr())
			c.disconnect <- c
			return
		}
		stringCmd := string(data)
		cmd, err := cftp.DeserializeCommand(strings.TrimSuffix(stringCmd, string(cftp.END_OF_MSG)))
		if err != nil {
			log.Printf("Error deserializing message: %v", err)
		}
		cmdChn <- cmd

		if cmd.Method == CMD_SEND {
			for i := 0; i < int(cmd.FileInfo.Size); {
				contentData, err := c.readFileContent(DEFAULT_BUFFER_SIZE, reader)
				if err != nil {
					log.Printf("Client %v disconected", c.conn.RemoteAddr())
					c.disconnect <- c
					return
				}
				i += len(contentData)
				contentChn <- contentData
			}
		}
	}
}

func (c *client) startWriter() {
	for {
		msg := <-c.writeChan
		_, err := c.conn.Write(msg)
		if err != nil {
			log.Println(err)
			break
		}
		//TODO: Tell the server that the connection is closed.
	}

}

func (c *client) readFileContent(bufSize int, reader *bufio.Reader) (data []byte, err error) {
	data = make([]byte, bufSize)
	n, err := reader.Read(data)
	data = data[:n]
	return
}
