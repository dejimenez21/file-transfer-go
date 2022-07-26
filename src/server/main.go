package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Hello! I'm the Server!")
	consoleOutputChannel := make(chan string)
	go handleConsoleOuput(consoleOutputChannel)
	listen(consoleOutputChannel)
}

func listen(ch chan<- string) {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, ch)
	}
}

func handleConnection(conn net.Conn, ch chan<- string) {
	log.SetPrefix("")
	ch <- fmt.Sprintln("Client connected from", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	data := make([]byte, 1000)
	for {
		n, err := reader.Read(data)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		data = data[:n]
		file, err := os.OpenFile("../../tests/scenarios/one-to-one-transfer/receiverFolder/test1.docx", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}

	// msg := fmt.Sprintf("Bytes received: %q", data)
	// ch <- msg
}

func handleConsoleOuput(ch chan string) {
	for {
		msg := <-ch
		fmt.Println(msg)
	}

}
