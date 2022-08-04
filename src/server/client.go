package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type client struct {
	conn      net.Conn
	writeChan chan []byte
	// cmdChan chan command
}

func (c *client) readRequest(cmdChn chan command, contentChn chan<- []byte) (req command) {
	for {
		reader := bufio.NewReader(c.conn)
		data, err := reader.ReadBytes(EOT)
		if err != nil {
			log.Printf("Error reading message from: %v", c.conn.RemoteAddr())
			return
		}
		stringCmd := string(data)
		cmd, err := deserializeCommand(strings.TrimSuffix(stringCmd, string(EOT)))
		if err != nil {
			log.Printf("Error deserializing message: %v", err)
		}
		cmdChn <- cmd

		if cmd.Method == CMD_SEND {
			for i := 0; i < int(cmd.FileInfo.Size); i += DEFAULT_BUFFER_SIZE {
				contentData, err := c.readFileContent(DEFAULT_BUFFER_SIZE, reader)
				if err != nil {
					//TODO: check EOF error
					log.Printf("an error occurred while reading file content from %s: %v", c.conn.RemoteAddr().String(), err)
				}
				contentChn <- contentData
			}
		}
	}
}

func (c *client) writeDelivery(message []byte) error {
	_, err := c.conn.Write(message)
	return err
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
