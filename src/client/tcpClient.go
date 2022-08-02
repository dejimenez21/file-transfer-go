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

func (c *tcpClient) readInput() (req request) {

	data, err := bufio.NewReader(c.conn).ReadBytes(EOT)
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading message from: %v", c.conn.RemoteAddr())
		}
		return
	}
	stringMsg := string(data)
	req, err = deserializeRequest(strings.TrimSuffix(stringMsg, string(EOT)))
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
	}
	return
	// chn <- &req
}

// filePath, err := getFilePath()
// if err != nil {
// 	log.Fatal(err)
// }
// conn, err := connect()
// if err != nil {
// 	os.Exit(1)
// }
// defer conn.Close()
// fileData, err := getFileBytes(filePath)
// if err != nil {
// 	os.Exit(1)
// }
// _, err = conn.Write(fileData)
// if err != nil {
// 	log.Fatal(err)
// }
