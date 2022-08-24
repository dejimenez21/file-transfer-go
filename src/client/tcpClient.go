package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type tcpClient struct {
	conn       net.Conn
	reader     *bufio.Reader
	serverAddr string
}

func (c *tcpClient) establishConnection() {
	if c.conn != nil {
		return
	}
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Couldn't establish connection: %v", err)
	}
	c.conn = conn
}

func (c *tcpClient) sendRequest(req request) error {
	c.establishConnection()
	serReq, err := serializeRequest(req)
	if err != nil {
		return fmt.Errorf("couldn't serialize request: %v", err)
	}
	serReq = append(serReq, EOT)
	_, err = c.conn.Write(serReq)
	if err != nil {
		return fmt.Errorf("couldn't send request: %v", err)
	}
	return nil
}

func (c *tcpClient) sendFileContent(content []byte) error {
	_, err := c.conn.Write(content)
	if err != nil {
		return fmt.Errorf("couldn't send request: %v", err)
	}
	return nil
}

func (c *tcpClient) readInput() (msg cftpMessage, err error) {
	if c.reader == nil {
		c.reader = bufio.NewReader(c.conn)
	}

	data, err := c.reader.ReadBytes(EOT)
	if err != nil {
		log.Fatalf("the connection with server was lost: %v", err)
	}
	stringMsg := string(data)

	if typeIndicator := (strings.SplitN(stringMsg, "\n", 2))[0]; typeIndicator == "chunk" {
		del, err := deserializeDelivery(stringMsg)
		if err != nil {
			return msg, err
		}
		del, err = c.readChunk(del)
		return &del, err
	}

	req, err := deserializeRequest(strings.TrimSuffix(stringMsg, string(EOT)))
	if err != nil {
		return msg, fmt.Errorf("error deserializing message: %v", err)
	}
	msg = &req
	return
}

func (c *tcpClient) readChunk(del delivery) (result delivery, err error) {
	data := make([]byte, del.Size)
	n, err := c.reader.Read(data)
	if err != nil {
		return result, fmt.Errorf("the connection with server was lost: %v", err)
	}
	data = data[:n]
	missing := del.Size - n
	for i := 0; i < missing; {
		missingData := make([]byte, missing-i)
		nm, err := c.conn.Read(missingData)
		if err != nil {
			return result, fmt.Errorf("the connection with server was lost: %v", err)
		}
		i += nm
		missingData = missingData[:nm]
		data = append(data, missingData...)
	}

	result = del
	result.Content = data
	return
}
