package main

import (
	"bufio"
	"io"
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

func (c *tcpClient) sendRequest(req request) {
	c.establishConnection()
	serReq, err := serializeRequest(req)
	if err != nil {
		log.Fatalf("Couldn't serialize request: %v", err)
	}
	serReq = append(serReq, EOT)
	_, err = c.conn.Write(serReq)
	if err != nil {
		log.Fatalf("Couldn't send request: %v", err)
	}
}

func (c *tcpClient) sendFileContent(content []byte) {
	_, err := c.conn.Write(content)
	if err != nil {
		log.Fatalf("Couldn't send request: %v", err)
	}
}

func (c *tcpClient) readInput() (msg cftpMessage, err error) {
	if c.reader == nil {
		c.reader = bufio.NewReader(c.conn)
	}
	// typeIndicator, err := reader.ReadString('\n')
	// if typeIndicator == "chunk\n" {

	// 	del := deserializeDelivery()
	// }
	data, err := c.reader.ReadBytes(EOT)
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading message from: %v", c.conn.RemoteAddr())
		}
		return
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
		log.Printf("Error deserializing message: %v", err)
	}
	msg = &req
	return
	// chn <- &req
}

func (c *tcpClient) readChunk(del delivery) (result delivery, err error) {
	data := make([]byte, del.Size)
	n, err := c.reader.Read(data)
	if err != nil {
		log.Fatal("lost connection with the server")
	}
	data = data[:n]
	if missing := del.Size - n; missing > 0 {
		missingData := make([]byte, missing)
		_, err := c.conn.Read(missingData)
		if err != nil {
			log.Fatal("lost connection with the server")
		}
		data = append(data, missingData...)
	}

	result = del
	result.Content = data
	return
}
