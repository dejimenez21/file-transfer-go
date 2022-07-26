package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	filePath, err := getFilePath()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := connect()
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()
	fileData, err := getFileBytes(filePath)
	if err != nil {
		os.Exit(1)
	}
	_, err = conn.Write(fileData)
	if err != nil {
		log.Fatal(err)
	}
}

func getFilePath() (path string, err error) {
	args := os.Args
	path = args[1]
	fmt.Println(path)
	return
}

func connect() (conn net.Conn, err error) {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err = d.DialContext(ctx, "tcp", "localhost:8888")
	if err != nil {
		log.Fatal(err)
	}

	return
}

func getFileBytes(path string) ([]byte, error) {
	fmt.Println("Getting file", path, "...")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	fileStat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileStat.Size()
	data := make([]byte, fileSize)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("read %d bytes: %q\n", count, data[:count])
	return data, err
}
