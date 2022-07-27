package main

import (
	"bufio"
	"log"
	"net"
)

type client struct {
	conn net.Conn
}

func (c *client) readCommand() {
	data, err := bufio.NewReader(c.conn).ReadBytes(0x04)
	if err != nil {
		log.Printf("Error reading command from: %v", c.conn.RemoteAddr())
		return
	}
	c.conn.Write(data)
	// for {
	// 	n, err := reader.Read(data)
	// 	if err != nil && err != io.EOF {
	// 		log.Fatal(err)
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}
	// 	data = data[:n]
	// 	file, err := os.OpenFile("../../tests/scenarios/one-to-one-transfer/receiverFolder/test1.docx", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()
	// 	_, err = file.Write(data)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
