package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type tcpClient struct {
	conn net.Conn
}

func (c *tcpClient) establishConnection() {
	if c.conn != nil {
		return
	}
	conn, err := net.Dial("tcp", "localhost:8888")
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
	reader := bufio.NewReader(c.conn)
	// typeIndicator, err := reader.ReadString('\n')
	// if typeIndicator == "chunk\n" {

	// 	del := deserializeDelivery()
	// }
	data, err := reader.ReadBytes(EOT)
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
		del, err = c.readChunk(reader, del)
		return del, err
	}

	msg, err = deserializeRequest(strings.TrimSuffix(stringMsg, string(EOT)))
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
	}
	return
	// chn <- &req
}

func (c *tcpClient) readChunk(reader *bufio.Reader, del delivery) (result delivery, err error) {
	data := make([]byte, del.Size)
	n, err := reader.Read(data)
	if err != nil {
		log.Fatal("lost connection with the server")
	}
	data = data[:n]
	result = del
	result.Content = data
	return
}
